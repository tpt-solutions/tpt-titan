// Speech Service for TTS/STT operations
class SpeechService {
    constructor() {
        this.baseURL = '/api/v1';
    }

    // Text-to-Speech operations
    async textToSpeech(text, modelId, options = {}) {
        try {
            const response = await fetch(`${this.baseURL}/speech/tts`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    text,
                    model_id: modelId,
                    options: {
                        voice: options.voice || 'alloy',
                        language: options.language || 'en',
                        speed: options.speed || 1.0,
                        pitch: options.pitch || 1.0,
                        volume: options.volume || 1.0,
                        audio_format: options.audioFormat || 'mp3'
                    }
                })
            });

            if (!response.ok) {
                throw new Error(`TTS request failed: ${response.statusText}`);
            }

            const result = await response.json();
            return result;
        } catch (error) {
            console.error('TTS error:', error);
            throw error;
        }
    }

    // Speech-to-Text operations
    async speechToText(audioData, modelId, options = {}) {
        try {
            const formData = new FormData();
            formData.append('audio', new Blob([audioData], { type: 'audio/wav' }));
            formData.append('model_id', modelId);
            formData.append('options', JSON.stringify({
                language: options.language || 'en',
                audio_format: options.audioFormat || 'wav'
            }));

            const response = await fetch(`${this.baseURL}/speech/stt`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error(`STT request failed: ${response.statusText}`);
            }

            const result = await response.json();
            return result;
        } catch (error) {
            console.error('STT error:', error);
            throw error;
        }
    }

    // Get available speech models
    async getAvailableModels(modelType = 'tts') {
        try {
            const response = await fetch(`${this.baseURL}/speech/models?type=${modelType}`);

            if (!response.ok) {
                throw new Error(`Failed to get models: ${response.statusText}`);
            }

            const result = await response.json();
            return result.models || [];
        } catch (error) {
            console.error('Get models error:', error);
            throw error;
        }
    }

    // Get speech request status
    async getRequestStatus(requestId) {
        try {
            const response = await fetch(`${this.baseURL}/speech/requests/${requestId}`);

            if (!response.ok) {
                throw new Error(`Failed to get status: ${response.statusText}`);
            }

            const result = await response.json();
            return result;
        } catch (error) {
            console.error('Get status error:', error);
            throw error;
        }
    }

    // Get user's speech settings
    async getSpeechSettings() {
        try {
            const response = await fetch(`${this.baseURL}/speech/settings`);

            if (!response.ok) {
                throw new Error(`Failed to get settings: ${response.statusText}`);
            }

            const result = await response.json();
            return result.settings;
        } catch (error) {
            console.error('Get settings error:', error);
            throw error;
        }
    }

    // Update user's speech settings
    async updateSpeechSettings(settings) {
        try {
            const response = await fetch(`${this.baseURL}/speech/settings`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ settings })
            });

            if (!response.ok) {
                throw new Error(`Failed to update settings: ${response.statusText}`);
            }

            const result = await response.json();
            return result;
        } catch (error) {
            console.error('Update settings error:', error);
            throw error;
        }
    }

    // Voice recording utilities
    startRecording(onDataAvailable, onError) {
        try {
            navigator.mediaDevices.getUserMedia({ audio: true })
                .then(stream => {
                    const mediaRecorder = new MediaRecorder(stream, {
                        mimeType: 'audio/webm;codecs=opus'
                    });

                    const audioChunks = [];

                    mediaRecorder.ondataavailable = (event) => {
                        audioChunks.push(event.data);
                    };

                    mediaRecorder.onstop = () => {
                        const audioBlob = new Blob(audioChunks, { type: 'audio/webm' });
                        onDataAvailable(audioBlob);
                    };

                    mediaRecorder.onerror = (error) => {
                        onError(error);
                    };

                    mediaRecorder.start();
                    return mediaRecorder;
                })
                .catch(error => {
                    onError(error);
                });
        } catch (error) {
            onError(error);
        }
    }

    // Convert audio blob to WAV format (basic implementation)
    async convertToWav(audioBlob) {
        return new Promise((resolve, reject) => {
            const audioContext = new (window.AudioContext || window.webkitAudioContext)();
            const fileReader = new FileReader();

            fileReader.onload = async (event) => {
                try {
                    const arrayBuffer = event.target.result;
                    const audioBuffer = await audioContext.decodeAudioData(arrayBuffer);

                    // Simple WAV conversion (for basic use cases)
                    // In production, you'd want a more robust WAV encoder
                    const wavBlob = await this.encodeWAV(audioBuffer);
                    resolve(wavBlob);
                } catch (error) {
                    reject(error);
                }
            };

            fileReader.onerror = reject;
            fileReader.readAsArrayBuffer(audioBlob);
        });
    }

    // Basic WAV encoder (simplified)
    async encodeWAV(audioBuffer) {
        const length = audioBuffer.length * audioBuffer.numberOfChannels * 2 + 44;
        const arrayBuffer = new ArrayBuffer(length);
        const view = new DataView(arrayBuffer);

        // WAV header
        const writeString = (offset, string) => {
            for (let i = 0; i < string.length; i++) {
                view.setUint8(offset + i, string.charCodeAt(i));
            }
        };

        writeString(0, 'RIFF');
        view.setUint32(4, length - 8, true);
        writeString(8, 'WAVE');
        writeString(12, 'fmt ');
        view.setUint32(16, 16, true);
        view.setUint16(20, 1, true);
        view.setUint16(22, audioBuffer.numberOfChannels, true);
        view.setUint32(24, audioBuffer.sampleRate, true);
        view.setUint32(28, audioBuffer.sampleRate * audioBuffer.numberOfChannels * 2, true);
        view.setUint16(32, audioBuffer.numberOfChannels * 2, true);
        view.setUint16(34, 16, true);
        writeString(36, 'data');
        view.setUint32(40, length - 44, true);

        // Audio data
        let offset = 44;
        for (let i = 0; i < audioBuffer.length; i++) {
            for (let channel = 0; channel < audioBuffer.numberOfChannels; channel++) {
                const sample = Math.max(-1, Math.min(1, audioBuffer.getChannelData(channel)[i]));
                view.setInt16(offset, sample * 0x7FFF, true);
                offset += 2;
            }
        }

        return new Blob([arrayBuffer], { type: 'audio/wav' });
    }
}

export default new SpeechService();
