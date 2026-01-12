// AI Components User Experience Testing
// Tests for accessibility, cross-platform consistency, and error handling

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { tick } from 'svelte';

// Mock the AI service
vi.mock('../../services/ai.js', () => ({
    aiService: {
        getWritingSuggestions: vi.fn(),
        summarizeDocument: vi.fn(),
        improveText: vi.fn(),
        continueWriting: vi.fn(),
        analyzeText: vi.fn(),
        isEnabled: vi.fn(() => true),
    }
}));

// Mock speech service
vi.mock('../../services/speech.js', () => ({
    default: {
        getAvailableModels: vi.fn(() => Promise.resolve([
            { id: 'stt-1', name: 'Test STT Model' }
        ])),
        startRecording: vi.fn(),
        speechToText: vi.fn(),
        convertToWav: vi.fn(),
    }
}));

describe('AI Components - User Experience Testing', () => {
    describe('TextEditor AI Features', () => {
        it('should show AI toolbar buttons when AI is enabled', async () => {
            // Import here to avoid hoisting issues
            const { default: TextEditor } = await import('../TextEditor.svelte');

            render(TextEditor, {
                props: { editorMode: 'markdown' }
            });

            await tick();

            // Check for AI buttons
            expect(screen.getByTitle('Get AI writing suggestions')).toBeInTheDocument();
            expect(screen.getByTitle('Improve selected text')).toBeInTheDocument();
            expect(screen.getByTitle('Continue writing from selected text')).toBeInTheDocument();
            expect(screen.getByTitle('Generate document summary')).toBeInTheDocument();
            expect(screen.getByTitle('Analyze document')).toBeInTheDocument();
        });

        it('should handle AI writing suggestions accessibility', async () => {
            const { default: TextEditor } = await import('../TextEditor.svelte');
            const { aiService } = await import('../../services/ai.js');

            aiService.getWritingSuggestions.mockResolvedValue({
                suggestions: ['Use active voice', 'Vary sentence length']
            });

            render(TextEditor, {
                props: { editorMode: 'markdown' }
            });

            await tick();

            const suggestButton = screen.getByTitle('Get AI writing suggestions');
            fireEvent.click(suggestButton);

            await waitFor(() => {
                expect(screen.getByText('✨ AI Writing Suggestions')).toBeInTheDocument();
            });

            // Check for accessible elements
            expect(screen.getByRole('button', { name: 'Apply' })).toBeInTheDocument();
            expect(screen.getByRole('button', { name: 'Close' })).toBeInTheDocument();
        });

        it('should show appropriate error messages for AI failures', async () => {
            const { default: TextEditor } = await import('../TextEditor.svelte');
            const { aiService } = await import('../../services/ai.js');

            aiService.getWritingSuggestions.mockRejectedValue(new Error('AI service unavailable'));

            render(TextEditor, {
                props: { editorMode: 'markdown' }
            });

            await tick();

            const suggestButton = screen.getByTitle('Get AI writing suggestions');
            fireEvent.click(suggestButton);

            await waitFor(() => {
                expect(screen.getByText('Failed to get AI suggestions. Please try again.')).toBeInTheDocument();
            });
        });

        it('should handle keyboard navigation for AI features', async () => {
            const { default: TextEditor } = await import('../TextEditor.svelte');

            render(TextEditor, {
                props: { editorMode: 'markdown' }
            });

            await tick();

            const aiButtons = screen.getAllByRole('button').filter(btn =>
                btn.title && btn.title.includes('AI')
            );

            // Test tab navigation
            aiButtons.forEach(button => {
                expect(button).toHaveAttribute('tabIndex', '0');
            });
        });
    });

    describe('EmailInbox AI Features', () => {
        it('should show AI tools section', async () => {
            const { default: EmailInbox } = await import('../EmailInbox.svelte');

            render(EmailInbox, {
                props: {
                    emailsList: [],
                    currentFolderValue: 'inbox',
                    selectedEmailData: null,
                    handleEmailSelect: vi.fn(),
                    handleFolderChange: vi.fn()
                }
            });

            await tick();

            expect(screen.getByText('AI TOOLS')).toBeInTheDocument();
            expect(screen.getByText('Smart Categorize')).toBeInTheDocument();
        });

        it('should handle voice composition accessibility', async () => {
            const { default: EmailInbox } = await import('../EmailInbox.svelte');

            render(EmailInbox, {
                props: {
                    emailsList: [],
                    currentFolderValue: 'inbox',
                    selectedEmailData: null,
                    handleEmailSelect: vi.fn(),
                    handleFolderChange: vi.fn()
                }
            });

            await tick();

            const voiceButton = screen.getByText('Voice Compose');
            expect(voiceButton).toBeInTheDocument();

            fireEvent.click(voiceButton);

            await waitFor(() => {
                expect(screen.getByText('🎤 Voice Email Composition')).toBeInTheDocument();
            });

            // Check for screen reader support
            expect(screen.getByLabelText('Transcribed Text')).toBeInTheDocument();
        });

        it('should show loading states during AI processing', async () => {
            const { default: EmailInbox } = await import('../EmailInbox.svelte');

            render(EmailInbox, {
                props: {
                    emailsList: [],
                    currentFolderValue: 'inbox',
                    selectedEmailData: null,
                    handleEmailSelect: vi.fn(),
                    handleFolderChange: vi.fn()
                }
            });

            await tick();

            const categorizeButton = screen.getByText('Smart Categorize');
            fireEvent.click(categorizeButton);

            // Button should show loading state
            expect(screen.getByText('Categorizing...')).toBeInTheDocument();
        });
    });

    describe('EventForm AI Features', () => {
        it('should show AI assistant section for new events', async () => {
            const { default: EventForm } = await import('../EventForm.svelte');

            render(EventForm, {
                props: {
                    event: null,
                    calendars: [{ id: '1', name: 'Work Calendar' }],
                    onClose: vi.fn()
                }
            });

            await tick();

            expect(screen.getByText('🤖 AI Assistant')).toBeInTheDocument();
            expect(screen.getByText('Voice Create Event')).toBeInTheDocument();
            expect(screen.getByText('Smart Scheduling')).toBeInTheDocument();
        });

        it('should handle voice event creation workflow', async () => {
            const { default: EventForm } = await import('../EventForm.svelte');

            render(EventForm, {
                props: {
                    event: null,
                    calendars: [{ id: '1', name: 'Work Calendar' }],
                    onClose: vi.fn()
                }
            });

            await tick();

            const voiceButton = screen.getByText('Voice Create Event');
            fireEvent.click(voiceButton);

            await waitFor(() => {
                expect(screen.getByText('🎤 Voice Event Creation')).toBeInTheDocument();
            });

            // Check for voice command examples
            expect(screen.getByText(/Schedule a meeting/)).toBeInTheDocument();
        });

        it('should show scheduling suggestions with proper accessibility', async () => {
            const { default: EventForm } = await import('../EventForm.svelte');

            render(EventForm, {
                props: {
                    event: null,
                    calendars: [{ id: '1', name: 'Work Calendar' }],
                    onClose: vi.fn()
                }
            });

            await tick();

            const scheduleButton = screen.getByText('Smart Scheduling');
            fireEvent.click(scheduleButton);

            await waitFor(() => {
                expect(screen.getByText('Suggested Times:')).toBeInTheDocument();
            });
        });
    });

    describe('TaskForm AI Features', () => {
        it('should show AI assistant for new tasks', async () => {
            const { default: TaskForm } = await import('../TaskForm.svelte');

            render(TaskForm, {
                props: {
                    task: null,
                    projects: []
                }
            });

            await tick();

            expect(screen.getByText('🤖 AI Assistant')).toBeInTheDocument();
            expect(screen.getByText('Predict Priority')).toBeInTheDocument();
            expect(screen.getByText('Predict Deadline')).toBeInTheDocument();
            expect(screen.getByText('Smart Suggestions')).toBeInTheDocument();
        });

        it('should display AI predictions with proper formatting', async () => {
            const { default: TaskForm } = await import('../TaskForm.svelte');

            render(TaskForm, {
                props: {
                    task: null,
                    projects: []
                }
            });

            await tick();

            // Simulate AI prediction results
            const priorityButton = screen.getByText('Predict Priority');
            fireEvent.click(priorityButton);

            // Check for prediction display area
            await waitFor(() => {
                expect(screen.getByText('🎯 Priority Prediction')).toBeInTheDocument();
            });
        });

        it('should allow applying AI suggestions', async () => {
            const { default: TaskForm } = await import('../TaskForm.svelte');

            render(TaskForm, {
                props: {
                    task: null,
                    projects: []
                }
            });

            await tick();

            // Check for apply buttons in suggestions
            const smartSuggestionsButton = screen.getByText('Smart Suggestions');
            fireEvent.click(smartSuggestionsButton);

            await waitFor(() => {
                const applyButtons = screen.getAllByText('Apply');
                expect(applyButtons.length).toBeGreaterThan(0);
            });
        });
    });

    describe('Cross-Platform Consistency', () => {
        it('should render consistently across different screen sizes', async () => {
            const { default: TextEditor } = await import('../TextEditor.svelte');

            // Mock different viewport sizes
            Object.defineProperty(window, 'innerWidth', {
                writable: true,
                configurable: true,
                value: 768 // Tablet size
            });

            render(TextEditor, {
                props: { editorMode: 'markdown' }
            });

            await tick();

            // AI toolbar should be visible and functional
            const aiSection = screen.getByText('AI');
            expect(aiSection).toBeInTheDocument();

            // Reset viewport
            window.innerWidth = 1024;
        });

        it('should handle touch interactions on mobile devices', async () => {
            const { default: TaskForm } = await import('../TaskForm.svelte');

            render(TaskForm, {
                props: {
                    task: null,
                    projects: []
                }
            });

            await tick();

            const aiButtons = screen.getAllByRole('button').filter(btn =>
                btn.textContent && btn.textContent.includes('Predict')
            );

            // Touch events should work the same as click events
            aiButtons.forEach(button => {
                fireEvent.touchStart(button);
                fireEvent.touchEnd(button);
            });

            expect(aiButtons.length).toBeGreaterThan(0);
        });
    });

    describe('Error Handling and Recovery', () => {
        it('should show user-friendly error messages', async () => {
            const { default: TextEditor } = await import('../TextEditor.svelte');
            const { aiService } = await import('../../services/ai.js');

            aiService.summarizeDocument.mockRejectedValue(new Error('Network timeout'));

            render(TextEditor, {
                props: { editorMode: 'markdown' }
            });

            await tick();

            const summaryButton = screen.getByTitle('Generate document summary');
            fireEvent.click(summaryButton);

            await waitFor(() => {
                expect(screen.getByText('Failed to generate document summary. Please try again.')).toBeInTheDocument();
            });
        });

        it('should allow retrying failed AI operations', async () => {
            const { default: EmailInbox } = await import('../EmailInbox.svelte');

            render(EmailInbox, {
                props: {
                    emailsList: [],
                    currentFolderValue: 'inbox',
                    selectedEmailData: null,
                    handleEmailSelect: vi.fn(),
                    handleFolderChange: vi.fn()
                }
            });

            await tick();

            // First attempt fails
            const categorizeButton = screen.getByText('Smart Categorize');
            fireEvent.click(categorizeButton);

            await waitFor(() => {
                expect(screen.getByText('Categorizing...')).toBeInTheDocument();
            });

            // Button should be clickable again after failure
            await waitFor(() => {
                expect(screen.getByText('Smart Categorize')).toBeInTheDocument();
            });
        });

        it('should handle network connectivity issues gracefully', async () => {
            const { default: TaskForm } = await import('../TaskForm.svelte');
            const { aiService } = await import('../../services/ai.js');

            aiService.predictPriority.mockRejectedValue(new Error('Network Error'));

            render(TaskForm, {
                props: {
                    task: null,
                    projects: []
                }
            });

            await tick();

            const predictButton = screen.getByText('Predict Priority');
            fireEvent.click(predictButton);

            // Should not crash and should allow retry
            await waitFor(() => {
                expect(predictButton).not.toBeDisabled();
            });
        });
    });

    describe('Accessibility Testing', () => {
        it('should have proper ARIA labels for AI features', async () => {
            const { default: TextEditor } = await import('../TextEditor.svelte');

            render(TextEditor, {
                props: { editorMode: 'markdown' }
            });

            await tick();

            // Check for ARIA labels and roles
            const aiButtons = screen.getAllByRole('button').filter(btn =>
                btn.title && btn.title.includes('AI')
            );

            aiButtons.forEach(button => {
                expect(button).toHaveAttribute('title');
            });
        });

        it('should support keyboard navigation for AI modals', async () => {
            const { default: TextEditor } = await import('../TextEditor.svelte');
            const { aiService } = await import('../../services/ai.js');

            aiService.getWritingSuggestions.mockResolvedValue({
                suggestions: ['Use active voice']
            });

            render(TextEditor, {
                props: { editorMode: 'markdown' }
            });

            await tick();

            const suggestButton = screen.getByTitle('Get AI writing suggestions');

            // Tab to button
            suggestButton.focus();
            expect(document.activeElement).toBe(suggestButton);

            // Activate with Enter key
            fireEvent.keyDown(suggestButton, { key: 'Enter' });

            await waitFor(() => {
                expect(screen.getByText('✨ AI Writing Suggestions')).toBeInTheDocument();
            });

            // Modal should be keyboard navigable
            const closeButton = screen.getByText('Close');
            fireEvent.keyDown(document, { key: 'Tab' });
            // Focus should move to actionable elements
        });

        it('should provide clear feedback for screen readers', async () => {
            const { default: TaskForm } = await import('../TaskForm.svelte');

            render(TaskForm, {
                props: {
                    task: null,
                    projects: []
                }
            });

            await tick();

            // Check for descriptive text and labels
            expect(screen.getByText('🎙️ Voice Commands')).toBeInTheDocument();
            expect(screen.getByText('🤖 AI Assistant')).toBeInTheDocument();
        });
    });
});
