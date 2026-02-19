import { c as create_ssr_component, e as emailSearchQuery, a as add_attribute, b as escape, d as each } from "./calendar.js";
class SpeechService {
  constructor() {
    this.baseURL = "/api/v1";
  }
  // Text-to-Speech operations
  async textToSpeech(text, modelId, options = {}) {
    try {
      const response = await fetch(`${this.baseURL}/speech/tts`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          text,
          model_id: modelId,
          options: {
            voice: options.voice || "alloy",
            language: options.language || "en",
            speed: options.speed || 1,
            pitch: options.pitch || 1,
            volume: options.volume || 1,
            audio_format: options.audioFormat || "mp3"
          }
        })
      });
      if (!response.ok) {
        throw new Error(`TTS request failed: ${response.statusText}`);
      }
      const result = await response.json();
      return result;
    } catch (error) {
      console.error("TTS error:", error);
      throw error;
    }
  }
  // Speech-to-Text operations
  async speechToText(audioData, modelId, options = {}) {
    try {
      const formData = new FormData();
      formData.append("audio", new Blob([audioData], { type: "audio/wav" }));
      formData.append("model_id", modelId);
      formData.append("options", JSON.stringify({
        language: options.language || "en",
        audio_format: options.audioFormat || "wav"
      }));
      const response = await fetch(`${this.baseURL}/speech/stt`, {
        method: "POST",
        body: formData
      });
      if (!response.ok) {
        throw new Error(`STT request failed: ${response.statusText}`);
      }
      const result = await response.json();
      return result;
    } catch (error) {
      console.error("STT error:", error);
      throw error;
    }
  }
  // Get available speech models
  async getAvailableModels(modelType = "tts") {
    try {
      const response = await fetch(`${this.baseURL}/speech/models?type=${modelType}`);
      if (!response.ok) {
        throw new Error(`Failed to get models: ${response.statusText}`);
      }
      const result = await response.json();
      return result.models || [];
    } catch (error) {
      console.error("Get models error:", error);
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
      console.error("Get status error:", error);
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
      console.error("Get settings error:", error);
      throw error;
    }
  }
  // Update user's speech settings
  async updateSpeechSettings(settings) {
    try {
      const response = await fetch(`${this.baseURL}/speech/settings`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ settings })
      });
      if (!response.ok) {
        throw new Error(`Failed to update settings: ${response.statusText}`);
      }
      const result = await response.json();
      return result;
    } catch (error) {
      console.error("Update settings error:", error);
      throw error;
    }
  }
  // Voice recording utilities
  startRecording(onDataAvailable, onError) {
    try {
      navigator.mediaDevices.getUserMedia({ audio: true }).then((stream) => {
        const mediaRecorder = new MediaRecorder(stream, {
          mimeType: "audio/webm;codecs=opus"
        });
        const audioChunks = [];
        mediaRecorder.ondataavailable = (event) => {
          audioChunks.push(event.data);
        };
        mediaRecorder.onstop = () => {
          const audioBlob = new Blob(audioChunks, { type: "audio/webm" });
          onDataAvailable(audioBlob);
        };
        mediaRecorder.onerror = (error) => {
          onError(error);
        };
        mediaRecorder.start();
        return mediaRecorder;
      }).catch((error) => {
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
    const writeString = (offset2, string) => {
      for (let i = 0; i < string.length; i++) {
        view.setUint8(offset2 + i, string.charCodeAt(i));
      }
    };
    writeString(0, "RIFF");
    view.setUint32(4, length - 8, true);
    writeString(8, "WAVE");
    writeString(12, "fmt ");
    view.setUint32(16, 16, true);
    view.setUint16(20, 1, true);
    view.setUint16(22, audioBuffer.numberOfChannels, true);
    view.setUint32(24, audioBuffer.sampleRate, true);
    view.setUint32(28, audioBuffer.sampleRate * audioBuffer.numberOfChannels * 2, true);
    view.setUint16(32, audioBuffer.numberOfChannels * 2, true);
    view.setUint16(34, 16, true);
    writeString(36, "data");
    view.setUint32(40, length - 44, true);
    let offset = 44;
    for (let i = 0; i < audioBuffer.length; i++) {
      for (let channel = 0; channel < audioBuffer.numberOfChannels; channel++) {
        const sample = Math.max(-1, Math.min(1, audioBuffer.getChannelData(channel)[i]));
        view.setInt16(offset, sample * 32767, true);
        offset += 2;
      }
    }
    return new Blob([arrayBuffer], { type: "audio/wav" });
  }
}
const SpeechService$1 = new SpeechService();
function formatDate$1(dateString) {
  if (!dateString) return "";
  const date = new Date(dateString);
  const now = /* @__PURE__ */ new Date();
  const diffTime = Math.abs(now - date);
  const diffDays = Math.ceil(diffTime / (1e3 * 60 * 60 * 24));
  if (diffDays === 1) {
    return "Today";
  } else if (diffDays === 2) {
    return "Yesterday";
  } else if (diffDays <= 7) {
    return date.toLocaleDateString("en-US", { weekday: "short" });
  } else {
    return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
  }
}
function truncateText(text, maxLength = 100) {
  if (!text) return "";
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + "...";
}
const EmailInbox = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let inboxCount;
  let starredCount;
  let sentCount;
  let { emailsList = [] } = $$props;
  let { currentFolderValue = "inbox" } = $$props;
  let { selectedEmailData = null } = $$props;
  let { handleEmailSelect } = $$props;
  let { handleFolderChange } = $$props;
  let searchQuery = "";
  let filteredEmails = [];
  let availableSTTModels = [];
  let selectedSTTModel = null;
  emailSearchQuery.subscribe((value) => {
    searchQuery = value;
    filterEmails();
  });
  onMount(async () => {
    try {
      availableSTTModels = await SpeechService$1.getAvailableModels("stt");
      if (availableSTTModels.length > 0) {
        selectedSTTModel = availableSTTModels[0];
      }
    } catch (error) {
      console.error("Failed to initialize AI features for email:", error);
    }
  });
  function filterEmails() {
    if (!searchQuery.trim()) {
      filteredEmails = emailsList;
    } else {
      const query = searchQuery.toLowerCase();
      filteredEmails = emailsList.filter((email) => email.subject?.toLowerCase().includes(query) || email.sender_name?.toLowerCase().includes(query) || email.sender_email?.toLowerCase().includes(query) || email.content?.toLowerCase().includes(query));
    }
  }
  if ($$props.emailsList === void 0 && $$bindings.emailsList && emailsList !== void 0) $$bindings.emailsList(emailsList);
  if ($$props.currentFolderValue === void 0 && $$bindings.currentFolderValue && currentFolderValue !== void 0) $$bindings.currentFolderValue(currentFolderValue);
  if ($$props.selectedEmailData === void 0 && $$bindings.selectedEmailData && selectedEmailData !== void 0) $$bindings.selectedEmailData(selectedEmailData);
  if ($$props.handleEmailSelect === void 0 && $$bindings.handleEmailSelect && handleEmailSelect !== void 0) $$bindings.handleEmailSelect(handleEmailSelect);
  if ($$props.handleFolderChange === void 0 && $$bindings.handleFolderChange && handleFolderChange !== void 0) $$bindings.handleFolderChange(handleFolderChange);
  {
    if (emailsList) {
      filterEmails();
    }
  }
  inboxCount = emailsList.filter((email) => email.folder === "inbox" && !email.is_read).length;
  starredCount = emailsList.filter((email) => email.is_starred).length;
  sentCount = emailsList.filter((email) => email.folder === "sent").length;
  return `<div class="flex flex-col h-full"> <div class="p-4 border-b border-gray-200 dark:border-gray-700"><div class="relative"><label for="email-search" class="sr-only" data-svelte-h="svelte-1c4t5t7">Search emails</label> <input id="email-search" type="text" placeholder="Search emails..." class="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent" aria-label="Search emails by subject, sender, or content"${add_attribute("value", searchQuery, 0)}> <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none" aria-hidden="true" data-svelte-h="svelte-1rbna79"><svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg></div></div></div>  <div class="border-b border-gray-200 dark:border-gray-700"><div class="p-2 space-y-1"><button class="${"w-full flex items-center justify-between px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors " + escape(
    currentFolderValue === "inbox" ? "bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300" : "text-gray-700 dark:text-gray-300",
    true
  )}"><div class="flex items-center space-x-2" data-svelte-h="svelte-ojx9gc"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg> <span>Inbox</span></div> ${inboxCount > 0 ? `<span class="bg-blue-600 text-white text-xs px-2 py-1 rounded-full min-w-[20px] text-center">${escape(inboxCount)}</span>` : ``}</button> <button class="${"w-full flex items-center justify-between px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors " + escape(
    currentFolderValue === "starred" ? "bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300" : "text-gray-700 dark:text-gray-300",
    true
  )}"><div class="flex items-center space-x-2" data-svelte-h="svelte-r2h7zm"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path></svg> <span>Starred</span></div> ${starredCount > 0 ? `<span class="bg-yellow-600 text-white text-xs px-2 py-1 rounded-full min-w-[20px] text-center">${escape(starredCount)}</span>` : ``}</button> <button class="${"w-full flex items-center justify-between px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors " + escape(
    currentFolderValue === "sent" ? "bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300" : "text-gray-700 dark:text-gray-300",
    true
  )}"><div class="flex items-center space-x-2" data-svelte-h="svelte-1tia16o"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path></svg> <span>Sent</span></div> ${sentCount > 0 ? `<span class="bg-gray-600 text-white text-xs px-2 py-1 rounded-full min-w-[20px] text-center">${escape(sentCount)}</span>` : ``}</button>  <div class="pt-2 border-t border-gray-200 dark:border-gray-600"><div class="px-3 py-1" data-svelte-h="svelte-1h0dmxz"><span class="text-xs text-gray-500 font-medium">AI TOOLS</span></div> <button ${""} class="w-full flex items-center px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-300 disabled:opacity-50"><div class="flex items-center space-x-2"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path></svg> <span>${escape("Smart Categorize")}</span></div></button> ${availableSTTModels.length > 0 ? `<button ${""} class="w-full flex items-center px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-300 disabled:opacity-50"><div class="flex items-center space-x-2"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg> <span>${escape("Voice Compose")}</span></div></button>` : ``}</div></div>  ${``}  ${``}</div>  <div class="flex-1 overflow-y-auto">${filteredEmails.length === 0 ? `<div class="p-8 text-center"><svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg> <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white" data-svelte-h="svelte-1549x6k">No emails</h3> <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">${escape(searchQuery ? "No emails match your search." : `No emails in ${currentFolderValue}.`)}</p></div>` : `<div class="divide-y divide-gray-200 dark:divide-gray-700">${each(filteredEmails, (email) => {
    return `<button class="${"w-full text-left p-4 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors " + escape(
      selectedEmailData?.id === email.id ? "bg-blue-50 dark:bg-blue-900 border-r-4 border-blue-500" : "",
      true
    ) + " focus:outline-none focus:ring-2 focus:ring-blue-500"}" aria-label="${"Read email: " + escape(email.subject || "(No subject)", true) + " from " + escape(email.sender_name || email.sender_email, true) + escape(email.is_starred ? ", starred" : "", true) + escape(!email.is_read ? ", unread" : "", true)}" type="button"><div class="flex items-start justify-between"><div class="flex-1 min-w-0"><div class="flex items-center space-x-2 mb-1">${email.is_starred ? `<svg class="w-4 h-4 text-yellow-500 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24"><path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path></svg>` : ``} <h4 class="text-sm font-medium text-gray-900 dark:text-white truncate">${escape(email.sender_name || email.sender_email)}</h4> ${!email.is_read ? `<span class="inline-block w-2 h-2 bg-blue-600 rounded-full flex-shrink-0"></span>` : ``}</div> <h3 class="text-sm font-medium text-gray-900 dark:text-white truncate mb-1">${escape(email.subject || "(No subject)")}</h3> <p class="text-sm text-gray-600 dark:text-gray-400 truncate">${escape(truncateText(email.content, 120))} </p></div> <div class="flex-shrink-0 ml-4 text-xs text-gray-500 dark:text-gray-400">${escape(formatDate$1(email.received_at || email.sent_at))} </div></div> </button>`;
  })}</div>`}</div></div>`;
});
function formatDate(dateString) {
  if (!dateString) return "";
  const date = new Date(dateString);
  return date.toLocaleString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
}
function formatEmailList(emails) {
  if (!emails || emails.length === 0) return "";
  return emails.join(", ");
}
const EmailViewer = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { email = null } = $$props;
  let { onClose } = $$props;
  async function markAsRead(isRead) {
    if (!email) return;
    try {
      const response = await fetch(`/api/v1/emails/${email.id}/read`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        body: JSON.stringify({ is_read: isRead })
      });
      if (response.ok) {
        email.is_read = isRead;
      }
    } catch (error) {
      console.error("Failed to mark email as read:", error);
    }
  }
  if ($$props.email === void 0 && $$bindings.email && email !== void 0) $$bindings.email(email);
  if ($$props.onClose === void 0 && $$bindings.onClose && onClose !== void 0) $$bindings.onClose(onClose);
  {
    if (email && !email.is_read) {
      markAsRead(true);
    }
  }
  return `${email ? `<div class="flex flex-col h-full bg-white dark:bg-gray-800"> <div class="border-b border-gray-200 dark:border-gray-700 p-6"><div class="flex items-center justify-between mb-4"><div class="flex items-center space-x-4"><button class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" data-svelte-h="svelte-1eec2rs"><svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg></button> <h2 class="text-xl font-semibold text-gray-900 dark:text-white truncate">${escape(email.subject || "(No subject)")}</h2></div> <div class="flex items-center space-x-2"> <button class="p-2 text-gray-400 hover:text-yellow-500 transition-colors"${add_attribute("title", email.is_starred ? "Remove star" : "Star email", 0)}><svg class="${"w-5 h-5 " + escape(email.is_starred ? "text-yellow-500 fill-current" : "", true)}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path></svg></button>  <div class="relative"><select class="text-sm border border-gray-300 dark:border-gray-600 rounded px-3 py-1 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"><option value="" data-svelte-h="svelte-1w43ca">Move to...</option><option value="inbox" data-svelte-h="svelte-1032kdm">Inbox</option><option value="archive" data-svelte-h="svelte-qa85pa">Archive</option><option value="trash" data-svelte-h="svelte-67g50u">Trash</option></select></div></div></div>  <div class="space-y-2"><div class="flex items-center justify-between"><div class="flex items-center space-x-4"><div class="flex-shrink-0"><div class="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold">${escape((email.sender_name || email.sender_email).charAt(0).toUpperCase())}</div></div> <div><h3 class="text-lg font-medium text-gray-900 dark:text-white">${escape(email.sender_name || email.sender_email)}</h3> <p class="text-sm text-gray-600 dark:text-gray-400">${escape(email.sender_email)}</p></div></div> <div class="text-sm text-gray-500 dark:text-gray-400">${escape(formatDate(email.received_at || email.sent_at))}</div></div> ${email.recipient_emails && email.recipient_emails.length > 0 ? `<div><span class="text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-pbqiax">To:</span> <span class="text-sm text-gray-900 dark:text-white ml-2">${escape(formatEmailList(email.recipient_emails))}</span></div>` : ``} ${email.cc_emails && email.cc_emails.length > 0 ? `<div><span class="text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-xqrpew">CC:</span> <span class="text-sm text-gray-900 dark:text-white ml-2">${escape(formatEmailList(email.cc_emails))}</span></div>` : ``}</div></div>  <div class="flex-1 overflow-y-auto p-6"><div class="max-w-none">${email.html_content ? ` <div class="prose dark:prose-invert max-w-none"><!-- HTML_TAG_START -->${email.html_content}<!-- HTML_TAG_END --></div>` : `${email.content ? ` <div class="whitespace-pre-wrap text-gray-900 dark:text-white font-mono text-sm leading-relaxed">${escape(email.content)}</div>` : `<div class="text-gray-500 dark:text-gray-400 italic" data-svelte-h="svelte-15c2rfs">No content available</div>`}`}</div></div>  ${email.attendees && email.attendees.length > 0 ? `<div class="border-t border-gray-200 dark:border-gray-700 p-6"><h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3" data-svelte-h="svelte-pqh91y">Attendees</h4> <div class="flex flex-wrap gap-2">${each(email.attendees, (attendee) => {
    return `<div class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">${escape(attendee.name)} ${attendee.status !== "pending" ? `<span class="ml-2 text-xs">(${escape(attendee.status)})
								</span>` : ``} </div>`;
  })}</div></div>` : ``}</div>` : ``}`;
});
export {
  EmailInbox as E,
  EmailViewer as a
};
