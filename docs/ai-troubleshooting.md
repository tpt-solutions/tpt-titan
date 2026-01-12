# TPT Titan AI Troubleshooting Guide

## Common Issues and Solutions

This guide covers the most common issues users encounter with TPT Titan's AI features, along with step-by-step solutions.

## Quick Diagnosis

### Is AI Working?

**Test 1: Basic AI Functionality**
1. Open the text editor
2. Type a sentence
3. Click the "✨ Suggest" button
4. If you see suggestions, AI is working

**Test 2: Voice Features**
1. Click any microphone button
2. Allow microphone permissions
3. Speak a short phrase
4. If text appears, voice input is working

**Test 3: Provider Status**
1. Go to Settings → AI Settings
2. Check that providers show "Connected" status
3. Verify API keys are entered (for cloud providers)

## AI Features Not Working

### Symptoms
- AI buttons are disabled or grayed out
- Features don't respond when clicked
- "AI service unavailable" messages

### Solutions

#### 1. Enable AI Features
```
Settings → AI Settings → Enable AI Features: ON
```

#### 2. Configure AI Providers
**For Local AI (Ollama)**:
1. Install Ollama from https://ollama.com
2. Pull required models:
   ```bash
   ollama pull qwen2.5:7b-instruct
   ollama pull qwen2.5-vl:7b
   ```
3. Verify Ollama is running: `ollama serve`

**For Cloud AI**:
1. Get API keys from providers:
   - OpenRouter: https://openrouter.ai/
   - ElevenLabs: https://elevenlabs.io/
   - AssemblyAI: https://assemblyai.com/
2. Enter keys in Settings → AI Settings → API Keys

#### 3. Check Network Connection
- Cloud AI features require internet
- Local AI works offline
- Check firewall settings for Ollama port (11434)

#### 4. Restart Application
- Close TPT Titan completely
- Wait 10 seconds
- Restart the application
- Try AI features again

## Voice Input Issues

### Symptoms
- Microphone doesn't activate
- "Microphone access denied" error
- Poor transcription quality
- No response to voice commands

### Solutions

#### Microphone Permissions
**Chrome/Edge**:
1. Click lock icon in address bar
2. Set microphone to "Allow"
3. Refresh the page

**Firefox**:
1. Click microphone icon in address bar
2. Select "Allow"
3. Choose your microphone device

**System Permissions**:
- Windows: Settings → Privacy → Microphone
- macOS: System Preferences → Security & Privacy → Microphone
- Linux: Check browser permissions

#### Audio Quality Issues
**Poor Transcription**:
1. Use a quality microphone (USB recommended)
2. Speak clearly and at normal pace
3. Reduce background noise
4. Test in a quiet environment

**Echo or Feedback**:
1. Use headphones to prevent audio feedback
2. Lower microphone volume in system settings
3. Close other applications using the microphone

#### Voice Command Recognition
**Commands Not Working**:
1. Speak clearly and pause between words
2. Use exact command phrases shown in help
3. Check that voice features are enabled in settings
4. Try different microphone devices

## Performance Issues

### Symptoms
- AI features are slow to respond
- High CPU/memory usage
- Application becomes unresponsive
- Features time out

### Solutions

#### Optimize Local AI Performance
**Model Selection**:
1. Use smaller models for basic tasks:
   ```bash
   ollama pull qwen2.5:3b-instruct  # Instead of 7B
   ```
2. Reserve large models for complex tasks

**Hardware Acceleration**:
1. Enable GPU acceleration in AI Settings
2. Ensure graphics drivers are updated
3. Check CUDA installation (NVIDIA GPUs)

**Memory Management**:
1. Enable "Low Power Mode" in AI Settings
2. Close unnecessary applications
3. Restart computer to free memory

#### Cloud AI Optimization
**API Rate Limits**:
1. Upgrade to higher API tier
2. Implement request throttling
3. Use local AI during peak hours

**Network Issues**:
1. Check internet speed (minimum 10 Mbps recommended)
2. Use wired connection instead of WiFi
3. Configure proxy settings if applicable

## API and Authentication Errors

### Symptoms
- "Invalid API key" errors
- "Rate limit exceeded" messages
- Authentication failures

### Solutions

#### API Key Issues
**Invalid Key**:
1. Verify key in Settings → AI Settings → API Keys
2. Check for extra spaces or characters
3. Regenerate key from provider if needed
4. Ensure correct provider selection

**Expired Key**:
1. Check key expiration date with provider
2. Generate new key if expired
3. Update key in TPT Titan settings
4. Test with a simple AI request

#### Rate Limiting
**Solutions**:
1. Check your API usage dashboard
2. Upgrade to higher usage tier
3. Implement usage monitoring in settings
4. Use local AI as backup during limits

## Document Processing Issues

### Symptoms
- OCR fails on documents
- Poor text extraction quality
- Large documents cause timeouts

### Solutions

#### OCR Quality Issues
**Poor Text Recognition**:
1. Ensure document is well-lit and clear
2. Use higher quality scans (300+ DPI)
3. Supported formats: PDF, PNG, JPG, TIFF
4. Try different document orientations

**Large Document Handling**:
1. Split large documents into smaller chunks
2. Process one page at a time
3. Use cloud AI for large documents
4. Enable "Batch Processing" in settings

#### Document Upload Issues
**Upload Failures**:
1. Check file size limits (max 50MB)
2. Verify supported file types
3. Ensure stable internet connection
4. Try uploading smaller documents first

## Email AI Issues

### Symptoms
- Email categorization fails
- Voice composition doesn't work
- Task conversion errors

### Solutions

#### Email Categorization
**Poor Categorization**:
1. Ensure emails have clear subject lines
2. Check that AI features are enabled
3. Try categorizing smaller batches
4. Review and provide feedback on results

#### Voice Email Composition
**Transcription Issues**:
1. Speak clearly and pause between sentences
2. Use proper email formatting commands
3. Check microphone settings
4. Try shorter email compositions first

## Calendar AI Issues

### Symptoms
- Voice event creation fails
- Scheduling suggestions not appearing
- Meeting summaries not generating

### Solutions

#### Voice Event Creation
**Parsing Errors**:
1. Use clear, natural language
2. Include specific dates, times, and durations
3. Example: "Schedule team meeting tomorrow at 2 PM for 1 hour"
4. Check voice input settings

#### Smart Scheduling
**No Suggestions**:
1. Ensure calendar has existing events
2. Check work hour preferences in settings
3. Verify AI features are enabled
4. Try during business hours

## Task Management AI Issues

### Symptoms
- Priority prediction fails
- Deadline suggestions not appearing
- Voice task creation errors

### Solutions

#### AI Task Prioritization
**Inaccurate Predictions**:
1. Provide more detailed task descriptions
2. Include urgency keywords
3. Add relevant tags
4. Review and adjust predictions manually

#### Task Deadline Prediction
**Poor Suggestions**:
1. Build task completion history
2. Provide more context in descriptions
3. Include specific requirements
4. Adjust suggestions based on your preferences

## Workflow Automation Issues

### Symptoms
- Workflows don't trigger
- Actions fail to execute
- AI optimization suggestions missing

### Solutions

#### Workflow Triggers
**Not Firing**:
1. Check trigger conditions are met
2. Verify workflow is active
3. Test with simple triggers first
4. Check workflow execution logs

#### Action Failures
**Common Issues**:
1. Verify permissions for automated actions
2. Check API connections for external services
3. Review error logs in workflow details
4. Test actions manually first

## Advanced Troubleshooting

### System Resource Monitoring

**Check System Resources**:
1. Open Task Manager (Windows) or Activity Monitor (macOS)
2. Monitor CPU, memory, and disk usage
3. Close resource-intensive applications
4. Restart system if resources are critically low

**AI Resource Usage**:
1. Check AI Settings → Performance
2. Monitor memory usage graphs
3. Enable "Low Power Mode" if needed
4. Switch to cloud AI for resource-intensive tasks

### Network Diagnostics

**Test Connectivity**:
1. Ping AI service endpoints
2. Check DNS resolution
3. Test with different network connections
4. Verify firewall and proxy settings

**API Endpoint Testing**:
1. Use browser dev tools to test API calls
2. Check response times and error codes
3. Verify API key authentication
4. Test with simple requests first

### Logs and Debugging

**Enable Debug Logging**:
1. Settings → AI Settings → Advanced
2. Enable "Debug Logging"
3. Check application logs for errors
4. Share logs with support if needed

**Console Debugging**:
1. Open browser developer tools (F12)
2. Check Console tab for JavaScript errors
3. Monitor Network tab for failed requests
4. Test API endpoints directly

## Getting Help

### Self-Service Options

**In-App Help**:
- Click "?" icons throughout the application
- Access contextual help for specific features
- Use built-in troubleshooting wizards

**Documentation**:
- Setup Guide: `docs/ai-setup-guide.md`
- Feature Tutorial: `docs/ai-features-tutorial.md`
- API Reference: `docs/api-reference.md`

### Support Channels

**Community Support**:
- Forum: https://community.tpt-titan.com
- GitHub Issues: Report bugs and request features
- User Groups: Join local TPT Titan user groups

**Professional Support**:
- Support Portal: https://support.tpt-titan.com
- Live Chat: Available during business hours
- Phone Support: Enterprise customers
- Email: support@tpt-titan.com

### Escalation Process

1. **Try Self-Service**: Check docs and common solutions
2. **Community Help**: Search forums and ask questions
3. **Support Ticket**: Create detailed ticket with logs
4. **Live Support**: Use chat for immediate assistance
5. **Escalation**: Enterprise customers can request priority support

## Prevention and Best Practices

### Regular Maintenance
- Keep AI providers updated
- Monitor API usage and costs
- Regularly test AI features
- Update TPT Titan to latest version

### Performance Optimization
- Use appropriate AI models for your hardware
- Enable caching for frequently used features
- Schedule heavy AI tasks during off-peak hours
- Monitor system resources regularly

### Security Best Practices
- Rotate API keys regularly
- Use local AI for sensitive data
- Monitor access logs
- Keep backup configurations

---

## Quick Reference

### Most Common Fixes

| Issue | Quick Fix |
|-------|-----------|
| AI not working | Enable in Settings → AI Settings |
| Voice not working | Check microphone permissions |
| Slow performance | Switch to cloud AI or smaller models |
| API errors | Verify API keys in settings |
| Network issues | Check internet connection |

### Emergency Contacts

- **Critical Issues**: support@tpt-titan.com (24/7)
- **Billing Issues**: billing@tpt-titan.com
- **Security Issues**: security@tpt-titan.com

---

*This troubleshooting guide is regularly updated. Check for the latest version at docs.tpt-titan.com/troubleshooting*
