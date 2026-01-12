# TPT Titan AI Setup Guide

## Overview

TPT Titan includes powerful AI features designed specifically for small and medium enterprises. This guide will help you configure and optimize AI services for your business needs.

## Quick Start

### 1. Enable AI Features

1. Log into TPT Titan
2. Go to **Settings** → **AI Settings**
3. Toggle **"Enable AI Features"** to ON
4. Configure your preferred AI providers

### 2. Choose Your AI Providers

TPT Titan supports multiple AI providers for maximum flexibility:

#### Local AI (Recommended for Privacy)
- **Ollama**: Run AI models locally on your hardware
- **Privacy-focused**: No data leaves your network
- **Cost-effective**: No API costs after initial setup

#### Cloud AI (Recommended for Performance)
- **OpenRouter**: Access to multiple AI models through one API
- **ElevenLabs**: High-quality text-to-speech
- **AssemblyAI**: Advanced speech-to-text

## Detailed Setup Instructions

### Local AI Setup with Ollama

#### Windows Installation

1. **Download Ollama**:
   ```
   Download from: https://ollama.com/download/windows
   ```

2. **Install and Start Ollama**:
   ```bash
   # Ollama will start automatically after installation
   # Verify it's running:
   ollama --version
   ```

3. **Pull Required Models**:
   ```bash
   # Writing and analysis models
   ollama pull qwen2.5:7b-instruct
   ollama pull qwen2.5-coder:7b-instruct

   # Vision models for document processing
   ollama pull qwen2.5-vl:7b
   ollama pull llama3.2-vision

   # Speech models (if using local TTS)
   ollama pull llama3.1:8b
   ```

4. **Configure TPT Titan**:
   - Go to **Settings** → **AI Settings**
   - Set **Default Local AI** to **Ollama**
   - Ollama URL: `http://localhost:11434` (default)

#### macOS Installation

1. **Install via Homebrew**:
   ```bash
   brew install ollama
   ```

2. **Start Ollama Service**:
   ```bash
   brew services start ollama
   ```

3. **Pull Models** (same as Windows instructions above)

#### Linux Installation

1. **Download and Install**:
   ```bash
   curl -fsSL https://ollama.com/install.sh | sh
   ```

2. **Start Ollama**:
   ```bash
   systemctl start ollama  # systemd
   # or
   service ollama start   # sysvinit
   ```

### Cloud AI Setup

#### OpenRouter Configuration

1. **Get API Key**:
   - Visit: https://openrouter.ai/
   - Sign up for an account
   - Generate an API key

2. **Configure in TPT Titan**:
   - Go to **Settings** → **AI Settings** → **API Keys**
   - Enter your OpenRouter API key
   - Set as default provider for text generation

#### ElevenLabs (Text-to-Speech)

1. **Get API Key**:
   - Visit: https://elevenlabs.io/
   - Create account and get API key

2. **Configure in TPT Titan**:
   - Add to **Speech Settings**
   - Choose your preferred voice
   - Set as default TTS provider

#### AssemblyAI (Speech-to-Text)

1. **Get API Key**:
   - Visit: https://www.assemblyai.com/
   - Sign up and get API key

2. **Configure in TPT Titan**:
   - Add to **Speech Settings**
   - Set as default STT provider

## Hardware Requirements

### Minimum Requirements
- **RAM**: 8GB (16GB recommended)
- **Storage**: 10GB free space for AI models
- **CPU**: Quad-core processor (8-core recommended)
- **GPU**: Optional but recommended for faster processing

### Recommended Configurations

#### Small Business (1-10 users)
- **Local AI**: Ollama with 7B parameter models
- **Cloud AI**: OpenRouter for peak usage
- **Hardware**: 16GB RAM, SSD storage

#### Medium Business (10-50 users)
- **Local AI**: Ollama with 14B+ parameter models
- **Cloud AI**: Dedicated OpenRouter plan
- **Hardware**: 32GB RAM, dedicated GPU recommended

## Performance Optimization

### Local AI Optimization

1. **Model Selection**:
   - Use smaller models for basic tasks (3B-7B parameters)
   - Reserve larger models for complex analysis (14B+ parameters)

2. **Hardware Acceleration**:
   - Enable GPU acceleration if available
   - Use CUDA for NVIDIA GPUs
   - Use Metal for Apple Silicon

3. **Memory Management**:
   - Monitor RAM usage in AI Settings
   - Enable "Low Power Mode" for battery-powered devices

### Network Optimization

1. **Caching**:
   - Enable response caching in AI Settings
   - Set appropriate cache TTL (Time To Live)

2. **Batch Processing**:
   - Use batch processing for multiple documents
   - Schedule heavy AI tasks during off-peak hours

## Troubleshooting

### Common Issues

#### Ollama Connection Failed
```
Error: Connection refused on localhost:11434
```

**Solutions**:
1. Verify Ollama is running: `ollama serve`
2. Check firewall settings
3. Confirm correct URL in AI Settings

#### Out of Memory Errors
```
Error: CUDA out of memory
```

**Solutions**:
1. Use smaller models
2. Enable "Low Power Mode"
3. Close other memory-intensive applications
4. Consider cloud AI for heavy workloads

#### Slow Performance
**Solutions**:
1. Enable hardware acceleration
2. Use SSD storage for model files
3. Reduce concurrent AI requests
4. Switch to cloud AI for demanding tasks

### API Key Issues

#### Invalid API Key
```
Error: Authentication failed
```

**Solutions**:
1. Verify API key in settings
2. Check API key expiration
3. Ensure correct provider selection

#### Rate Limiting
```
Error: Rate limit exceeded
```

**Solutions**:
1. Upgrade to higher API tier
2. Implement request throttling
3. Use local AI during peak hours
4. Cache frequently used responses

## Security Considerations

### Local AI Security
- ✅ Data stays on your network
- ✅ No external API calls
- ✅ Full control over data processing

### Cloud AI Security
- 🔒 End-to-end encryption for API calls
- 🔒 SOC 2 compliant providers
- 🔒 Data minimization practices
- ⚠️ Review provider privacy policies

## Cost Optimization

### Local AI Costs
- **Setup Cost**: Free (one-time hardware investment)
- **Operating Cost**: Electricity only
- **Scaling**: Limited by hardware

### Cloud AI Costs
- **Pay-per-use**: Based on tokens/requests
- **Free Tiers**: Available for light usage
- **Enterprise Plans**: Volume discounts available

### Cost Monitoring
- View usage statistics in **Settings** → **AI Usage**
- Set budget alerts for cloud API usage
- Track token consumption by feature

## Best Practices

### For Small Businesses
1. Start with local AI for basic features
2. Use cloud AI for advanced capabilities
3. Monitor costs and usage patterns
4. Train team members on AI features

### For IT Administrators
1. Implement user access controls
2. Set up monitoring and alerting
3. Plan for scaling as usage grows
4. Maintain backup configurations

### For End Users
1. Start with simple AI features (writing assistance)
2. Gradually explore advanced features
3. Provide feedback for continuous improvement
4. Use voice features for hands-free operation

## Support and Resources

### Getting Help
- **Documentation**: https://docs.tpt-titan.com/ai
- **Community Forum**: https://community.tpt-titan.com
- **Support Portal**: https://support.tpt-titan.com

### Additional Resources
- **AI Best Practices Guide**: Learn optimization techniques
- **Video Tutorials**: Step-by-step feature walkthroughs
- **API Reference**: Technical integration details

---

*This guide is regularly updated. Check for the latest version at docs.tpt-titan.com/ai-setup*
