package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"tpt-titan/backend/config"
)

// P2PService handles peer-to-peer networking for local collaboration
type P2PService struct {
	config     *config.P2PConfig
	peers      map[string]*Peer // peerID -> Peer
	peerMutex  sync.RWMutex
	discovery  *PeerDiscovery
	relay      *CloudRelay
	listener   net.Listener
	running    bool
	stopChan   chan struct{}

	// Connection management
	localPeers   map[string]*Peer // Direct local network peers
	remotePeers  map[string]*Peer // Remote peers via relay
	peerID       string           // This instance's unique peer ID

	// Callbacks for handling spreadsheet events
	onPeerConnected    func(peerID string)
	onPeerDisconnected func(peerID string)
	onSpreadsheetSync  func(peerID string, spreadsheetID uuid.UUID, data map[string]interface{})
	onCellUpdate       func(peerID string, spreadsheetID uuid.UUID, cellRef string, value interface{})
}

// Peer represents a connected peer in the P2P network
type Peer struct {
	ID          string
	Address     string
	Name        string
	LastSeen    time.Time
	Connection  net.Conn
	Encoder     *json.Encoder
	Decoder     *json.Decoder
	IsConnected bool
	mutex       sync.Mutex
}

// P2PMessage represents messages exchanged between peers
type P2PMessage struct {
	Type         string                 `json:"type"`
	PeerID       string                 `json:"peer_id"`
	Timestamp    time.Time              `json:"timestamp"`
	SpreadsheetID uuid.UUID             `json:"spreadsheet_id,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
}

// PeerDiscovery handles mDNS peer discovery
type PeerDiscovery struct {
	serviceName string
	serviceType string
	port        int
	peers       map[string]string // peerID -> address
	mutex       sync.RWMutex
}

// CloudRelay handles remote connections via cloud relay servers
type CloudRelay struct {
	serverURL  string
	token      string // Authentication token
	connected  bool
	conn       net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
	mutex      sync.Mutex
}

// RelayMessage represents messages sent through the cloud relay
type RelayMessage struct {
	Type      string                 `json:"type"`
	Token     string                 `json:"token"`
	TargetPeer string                 `json:"target_peer,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewP2PService creates a new P2P service
func NewP2PService(cfg *config.P2PConfig) *P2PService {
	return &P2PService{
		config:     cfg,
		peers:      make(map[string]*Peer),
		stopChan:   make(chan struct{}),
		discovery:  NewPeerDiscovery(cfg.ServiceName, cfg.ServiceType, cfg.Port),
	}
}

// IsRunning reports whether the P2P service is currently active
func (p2p *P2PService) IsRunning() bool {
	return p2p.running
}

// Start initializes the P2P service
func (p2p *P2PService) Start() error {
	if p2p.running {
		return fmt.Errorf("P2P service already running")
	}

	log.Printf("Starting P2P service on port %d", p2p.config.Port)
	p2p.running = true

	// Start peer discovery
	if err := p2p.discovery.Start(); err != nil {
		return fmt.Errorf("failed to start peer discovery: %v", err)
	}

	// Start listening for connections
	if err := p2p.startListener(); err != nil {
		p2p.discovery.Stop()
		return fmt.Errorf("failed to start listener: %v", err)
	}

	// Start peer discovery loop
	go p2p.discoveryLoop()

	// Start connection management
	go p2p.connectionManager()

	log.Printf("P2P service started successfully")
	return nil
}

// Stop shuts down the P2P service
func (p2p *P2PService) Stop() {
	if !p2p.running {
		return
	}

	log.Printf("Stopping P2P service")
	p2p.running = false
	close(p2p.stopChan)

	// Stop discovery
	p2p.discovery.Stop()

	// Close listener
	if p2p.listener != nil {
		p2p.listener.Close()
	}

	// Disconnect all peers
	p2p.peerMutex.Lock()
	for peerID, peer := range p2p.peers {
		if peer.IsConnected {
			peer.Connection.Close()
		}
		delete(p2p.peers, peerID)
	}
	p2p.peerMutex.Unlock()

	log.Printf("P2P service stopped")
}

// startListener starts listening for incoming peer connections
func (p2p *P2PService) startListener() error {
	var listener net.Listener
	var err error

	if p2p.config.EnableEncryption {
		// TLS listener (would need certificates in production)
		cert, err := tls.X509KeyPair([]byte("dummy-cert"), []byte("dummy-key"))
		if err != nil {
			return fmt.Errorf("failed to create TLS certificate: %v", err)
		}

		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err = tls.Listen("tcp", fmt.Sprintf(":%d", p2p.config.Port), tlsConfig)
	} else {
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", p2p.config.Port))
	}

	if err != nil {
		return err
	}

	p2p.listener = listener

	// Accept connections in background
	go func() {
		for p2p.running {
			conn, err := listener.Accept()
			if err != nil {
				if p2p.running {
					log.Printf("Failed to accept connection: %v", err)
				}
				continue
			}

			go p2p.handleConnection(conn)
		}
	}()

	return nil
}

// handleConnection handles incoming peer connections
func (p2p *P2PService) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create decoder/encoder
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	// Read handshake message
	var handshake P2PMessage
	if err := decoder.Decode(&handshake); err != nil {
		log.Printf("Failed to read handshake: %v", err)
		return
	}

	if handshake.Type != "handshake" {
		log.Printf("Invalid handshake message type: %s", handshake.Type)
		return
	}

	peerID := handshake.PeerID
	peerName := handshake.Data["name"].(string)

	// Create peer
	peer := &Peer{
		ID:          peerID,
		Address:     conn.RemoteAddr().String(),
		Name:        peerName,
		LastSeen:    time.Now(),
		Connection:  conn,
		Encoder:     encoder,
		Decoder:     decoder,
		IsConnected: true,
	}

	// Add peer
	p2p.peerMutex.Lock()
	p2p.peers[peerID] = peer
	p2p.peerMutex.Unlock()

	log.Printf("Peer connected: %s (%s)", peerName, peerID)

	// Send handshake response
	response := P2PMessage{
		Type:      "handshake_ack",
		PeerID:    "self", // Would be actual peer ID
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"name": "TPT Titan Instance", // Would be configurable
		},
	}

	if err := encoder.Encode(response); err != nil {
		log.Printf("Failed to send handshake response: %v", err)
		return
	}

	// Notify about new peer
	if p2p.onPeerConnected != nil {
		p2p.onPeerConnected(peerID)
	}

	// Handle messages from this peer
	p2p.handlePeerMessages(peer)
}

// handlePeerMessages handles messages from a connected peer
func (p2p *P2PService) handlePeerMessages(peer *Peer) {
	for p2p.running && peer.IsConnected {
		var msg P2PMessage
		if err := peer.Decoder.Decode(&msg); err != nil {
			log.Printf("Failed to read message from peer %s: %v", peer.ID, err)
			break
		}

		p2p.handleMessage(peer.ID, msg)
	}

	// Peer disconnected
	p2p.peerMutex.Lock()
	if peer.IsConnected {
		peer.IsConnected = false
		peer.Connection.Close()
	}
	p2p.peerMutex.Unlock()

	if p2p.onPeerDisconnected != nil {
		p2p.onPeerDisconnected(peer.ID)
	}

	log.Printf("Peer disconnected: %s", peer.ID)
}

// handleMessage processes messages from peers
func (p2p *P2PService) handleMessage(peerID string, msg P2PMessage) {
	switch msg.Type {
	case "spreadsheet_sync":
		if p2p.onSpreadsheetSync != nil {
			spreadsheetID := msg.SpreadsheetID
			data := msg.Data
			p2p.onSpreadsheetSync(peerID, spreadsheetID, data)
		}

	case "cell_update":
		if p2p.onCellUpdate != nil {
			spreadsheetID := msg.SpreadsheetID
			cellRef := msg.Data["cell_ref"].(string)
			value := msg.Data["value"]
			p2p.onCellUpdate(peerID, spreadsheetID, cellRef, value)
		}

	case "ping":
		// Respond with pong
		p2p.sendMessage(peerID, P2PMessage{
			Type:      "pong",
			PeerID:    "self",
			Timestamp: time.Now(),
		})

	default:
		log.Printf("Unknown message type: %s from peer %s", msg.Type, peerID)
	}
}

// ConnectToPeer attempts to connect to a discovered peer
func (p2p *P2PService) ConnectToPeer(peerID, address string) error {
	p2p.peerMutex.Lock()
	if _, exists := p2p.peers[peerID]; exists {
		p2p.peerMutex.Unlock()
		return fmt.Errorf("already connected to peer %s", peerID)
	}
	p2p.peerMutex.Unlock()

	// Create connection
	var conn net.Conn
	var err error

	if p2p.config.EnableEncryption {
		// TLS connection (would need proper certificate validation)
		tlsConfig := &tls.Config{InsecureSkipVerify: true} // For local network
		conn, err = tls.Dial("tcp", address, tlsConfig)
	} else {
		conn, err = net.Dial("tcp", address)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to peer: %v", err)
	}

	// Send handshake
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	handshake := P2PMessage{
		Type:      "handshake",
		PeerID:    "self", // Would be actual peer ID
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"name": "TPT Titan Instance",
		},
	}

	if err := encoder.Encode(handshake); err != nil {
		conn.Close()
		return fmt.Errorf("failed to send handshake: %v", err)
	}

	// Wait for response
	var response P2PMessage
	if err := decoder.Decode(&response); err != nil {
		conn.Close()
		return fmt.Errorf("failed to read handshake response: %v", err)
	}

	if response.Type != "handshake_ack" {
		conn.Close()
		return fmt.Errorf("invalid handshake response: %s", response.Type)
	}

	peerName := response.Data["name"].(string)

	// Create peer
	peer := &Peer{
		ID:          peerID,
		Address:     address,
		Name:        peerName,
		LastSeen:    time.Now(),
		Connection:  conn,
		Encoder:     encoder,
		Decoder:     decoder,
		IsConnected: true,
	}

	// Add peer
	p2p.peerMutex.Lock()
	p2p.peers[peerID] = peer
	p2p.peerMutex.Unlock()

	log.Printf("Connected to peer: %s (%s)", peerName, peerID)

	// Notify about new peer
	if p2p.onPeerConnected != nil {
		p2p.onPeerConnected(peerID)
	}

	// Start message handler
	go p2p.handlePeerMessages(peer)

	return nil
}

// BroadcastMessage sends a message to all connected peers
func (p2p *P2PService) BroadcastMessage(msg P2PMessage) {
	p2p.peerMutex.RLock()
	peers := make([]*Peer, 0, len(p2p.peers))
	for _, peer := range p2p.peers {
		if peer.IsConnected {
			peers = append(peers, peer)
		}
	}
	p2p.peerMutex.RUnlock()

	for _, peer := range peers {
		if err := peer.Encoder.Encode(msg); err != nil {
			log.Printf("Failed to send message to peer %s: %v", peer.ID, err)
			peer.IsConnected = false
		}
	}
}

// sendMessage sends a message to a specific peer
func (p2p *P2PService) sendMessage(peerID string, msg P2PMessage) error {
	p2p.peerMutex.RLock()
	peer, exists := p2p.peers[peerID]
	p2p.peerMutex.RUnlock()

	if !exists || !peer.IsConnected {
		return fmt.Errorf("peer %s not connected", peerID)
	}

	peer.mutex.Lock()
	defer peer.mutex.Unlock()

	return peer.Encoder.Encode(msg)
}

// SyncSpreadsheet broadcasts spreadsheet data to peers
func (p2p *P2PService) SyncSpreadsheet(spreadsheetID uuid.UUID, data map[string]interface{}) {
	msg := P2PMessage{
		Type:          "spreadsheet_sync",
		PeerID:        "self",
		Timestamp:     time.Now(),
		SpreadsheetID: spreadsheetID,
		Data:          data,
	}

	p2p.BroadcastMessage(msg)
}

// UpdateCell broadcasts a cell update to peers
func (p2p *P2PService) UpdateCell(spreadsheetID uuid.UUID, cellRef string, value interface{}) {
	msg := P2PMessage{
		Type:          "cell_update",
		PeerID:        "self",
		Timestamp:     time.Now(),
		SpreadsheetID: spreadsheetID,
		Data: map[string]interface{}{
			"cell_ref": cellRef,
			"value":    value,
		},
	}

	p2p.BroadcastMessage(msg)
}

// GetConnectedPeers returns list of connected peers
func (p2p *P2PService) GetConnectedPeers() []map[string]interface{} {
	p2p.peerMutex.RLock()
	defer p2p.peerMutex.RUnlock()

	peers := make([]map[string]interface{}, 0, len(p2p.peers))
	for _, peer := range p2p.peers {
		if peer.IsConnected {
			peers = append(peers, map[string]interface{}{
				"id":       peer.ID,
				"name":     peer.Name,
				"address":  peer.Address,
				"last_seen": peer.LastSeen,
			})
		}
	}

	return peers
}

// discoveryLoop periodically discovers new peers
func (p2p *P2PService) discoveryLoop() {
	ticker := time.NewTicker(time.Duration(p2p.config.DiscoveryTimeout) * time.Second)
	defer ticker.Stop()

	for p2p.running {
		select {
		case <-ticker.C:
			p2p.discoverPeers()
		case <-p2p.stopChan:
			return
		}
	}
}

// discoverPeers discovers available peers on the network
func (p2p *P2PService) discoverPeers() {
	// This would implement mDNS/Bonjour discovery
	// For now, just log that discovery is running
	log.Printf("Running peer discovery...")

	// In a full implementation, this would:
	// 1. Send mDNS queries
	// 2. Listen for service advertisements
	// 3. Update peer list with discovered services
	// 4. Attempt connections to new peers
}

// connectionManager manages peer connections and health
func (p2p *P2PService) connectionManager() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for p2p.running {
		select {
		case <-ticker.C:
			p2p.checkPeerHealth()
		case <-p2p.stopChan:
			return
		}
	}
}

// checkPeerHealth checks if peers are still responsive
func (p2p *P2PService) checkPeerHealth() {
	p2p.peerMutex.Lock()
	defer p2p.peerMutex.Unlock()

	for peerID, peer := range p2p.peers {
		if peer.IsConnected {
			// Send ping
			pingMsg := P2PMessage{
				Type:      "ping",
				PeerID:    "self",
				Timestamp: time.Now(),
			}

			if err := peer.Encoder.Encode(pingMsg); err != nil {
				log.Printf("Peer %s appears disconnected: %v", peerID, err)
				peer.IsConnected = false
				peer.Connection.Close()

				if p2p.onPeerDisconnected != nil {
					p2p.onPeerDisconnected(peerID)
				}
			}
		}
	}
}

// SetCallbacks sets callback functions for handling events
func (p2p *P2PService) SetCallbacks(
	onPeerConnected func(peerID string),
	onPeerDisconnected func(peerID string),
	onSpreadsheetSync func(peerID string, spreadsheetID uuid.UUID, data map[string]interface{}),
	onCellUpdate func(peerID string, spreadsheetID uuid.UUID, cellRef string, value interface{}),
) {
	p2p.onPeerConnected = onPeerConnected
	p2p.onPeerDisconnected = onPeerDisconnected
	p2p.onSpreadsheetSync = onSpreadsheetSync
	p2p.onCellUpdate = onCellUpdate
}

// NewPeerDiscovery creates a new peer discovery service
func NewPeerDiscovery(serviceName, serviceType string, port int) *PeerDiscovery {
	return &PeerDiscovery{
		serviceName: serviceName,
		serviceType: serviceType,
		port:        port,
		peers:       make(map[string]string),
	}
}

// Start starts peer discovery
func (pd *PeerDiscovery) Start() error {
	// This would start mDNS service discovery
	// For now, just return success
	log.Printf("Starting peer discovery for service %s", pd.serviceName)
	return nil
}

// Stop stops peer discovery
func (pd *PeerDiscovery) Stop() {
	log.Printf("Stopping peer discovery")
}

// GetPeers returns discovered peers
func (pd *PeerDiscovery) GetPeers() map[string]string {
	pd.mutex.RLock()
	defer pd.mutex.RUnlock()

	peers := make(map[string]string)
	for k, v := range pd.peers {
		peers[k] = v
	}
	return peers
}
