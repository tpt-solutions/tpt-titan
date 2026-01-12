// AI Service for TPT Titan frontend
// Handles AI writing assistance, document summarization, and other AI features

const API_BASE = '/api/v1';

class AIService {
    constructor() {
        this.availableModels = [];
        this.userSettings = null;
    }

    // Initialize the service
    async initialize() {
        try {
            // Load available AI models
            await this.loadAvailableModels();

            // Load user AI settings
            await this.loadUserSettings();

            return true;
        } catch (error) {
            console.error('Failed to initialize AI service:', error);
            return false;
        }
    }

    // Load available AI models
    async loadAvailableModels() {
        try {
            const response = await fetch(`${API_BASE}/ai/models`, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (response.ok) {
                const data = await response.json();
                this.availableModels = data.models || [];
            } else {
                console.warn('Failed to load AI models');
                this.availableModels = [];
            }
        } catch (error) {
            console.error('Error loading AI models:', error);
            this.availableModels = [];
        }
    }

    // Load user AI settings
    async loadUserSettings() {
        try {
            const response = await fetch(`${API_BASE}/ai/settings`, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (response.ok) {
                this.userSettings = await response.json();
            } else {
                console.warn('Failed to load AI settings');
                this.userSettings = {
                    enable_ai_features: true,
                    enable_writing_assistance: true,
                    enable_document_summary: true,
                    default_writing_model: null
                };
            }
        } catch (error) {
            console.error('Error loading AI settings:', error);
            this.userSettings = {
                enable_ai_features: true,
                enable_writing_assistance: true,
                enable_document_summary: true,
                default_writing_model: null
            };
        }
    }

    // Get writing assistance suggestions
    async getWritingSuggestions(text, context = 'general') {
        if (!this.userSettings?.enable_writing_assistance) {
            return null;
        }

        try {
            const response = await fetch(`${API_BASE}/ai/writing-assistance`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    text: text,
                    context: context,
                    model_id: this.userSettings.default_writing_model
                })
            });

            if (response.ok) {
                return await response.json();
            } else {
                console.warn('Failed to get writing suggestions');
                return null;
            }
        } catch (error) {
            console.error('Error getting writing suggestions:', error);
            return null;
        }
    }

    // Generate document summary
    async summarizeDocument(content, summaryType = 'concise') {
        if (!this.userSettings?.enable_document_summary) {
            return null;
        }

        try {
            const response = await fetch(`${API_BASE}/ai/document-summary`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    content: content,
                    summary_type: summaryType,
                    model_id: this.userSettings.default_writing_model
                })
            });

            if (response.ok) {
                const result = await response.json();
                return result.summary;
            } else {
                console.warn('Failed to generate document summary');
                return null;
            }
        } catch (error) {
            console.error('Error generating document summary:', error);
            return null;
        }
    }

    // Improve text (grammar, style, clarity)
    async improveText(text, improvementType = 'general') {
        try {
            const response = await fetch(`${API_BASE}/ai/improve-text`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    text: text,
                    improvement_type: improvementType,
                    model_id: this.userSettings?.default_writing_model
                })
            });

            if (response.ok) {
                const result = await response.json();
                return result.improved_text;
            } else {
                console.warn('Failed to improve text');
                return text; // Return original if improvement fails
            }
        } catch (error) {
            console.error('Error improving text:', error);
            return text;
        }
    }

    // Continue writing (generate next sentences/paragraphs)
    async continueWriting(text, style = 'matching', length = 'medium') {
        try {
            const response = await fetch(`${API_BASE}/ai/continue-writing`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    text: text,
                    style: style,
                    length: length,
                    model_id: this.userSettings?.default_writing_model
                })
            });

            if (response.ok) {
                const result = await response.json();
                return result.continuation;
            } else {
                console.warn('Failed to continue writing');
                return '';
            }
        } catch (error) {
            console.error('Error continuing writing:', error);
            return '';
        }
    }

    // Analyze text for various metrics
    async analyzeText(text) {
        try {
            const response = await fetch(`${API_BASE}/ai/analyze-text`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    text: text
                })
            });

            if (response.ok) {
                return await response.json();
            } else {
                console.warn('Failed to analyze text');
                return null;
            }
        } catch (error) {
            console.error('Error analyzing text:', error);
            return null;
        }
    }

    // Check if AI features are enabled
    isEnabled() {
        return this.userSettings?.enable_ai_features === true;
    }

    // Get available writing models
    getWritingModels() {
        return this.availableModels.filter(model =>
            model.capabilities?.includes('writing') ||
            model.capabilities?.includes('text_generation')
        );
    }

    // Get user settings
    getSettings() {
        return this.userSettings;
    }

    // Update user settings
    async updateSettings(settings) {
        try {
            const response = await fetch(`${API_BASE}/ai/settings`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify(settings)
            });

            if (response.ok) {
                this.userSettings = await response.json();
                return true;
            } else {
                console.warn('Failed to update AI settings');
                return false;
            }
        } catch (error) {
            console.error('Error updating AI settings:', error);
            return false;
        }
    }
}

// Create and export a singleton instance
export const aiService = new AIService();
export default aiService;
