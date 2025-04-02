<template>
  <div class="w-80 flex flex-col overflow-hidden">
    <!-- Sidebar Header -->
    <RoomSidebarHeader
      :is-chat-open="isChatOpen"
      @close-sidebar="$emit('close-sidebar')"
    />

    <!-- Chat Content -->
    <RoomSidebarChat v-if="isChatOpen" :chat-messages="chatMessages" />

    <!-- Participants Content -->
    <RoomSidebarParticipants
      v-if="isParticipantsOpen"
      :participants="participants"
    />

    <!-- Chat Input -->
    <RoomChatInput v-if="isChatOpen" />
  </div>
</template>

<script setup lang="ts">
import {
  RoomChatInput,
  RoomSidebarChat,
  RoomSidebarHeader,
  RoomSidebarParticipants,
} from "#components";
import type { ChatMessage, Participant } from "~/types";

defineProps({
  isChatOpen: Boolean,
  isParticipantsOpen: Boolean,
  chatMessages: {
    type: Array as () => ChatMessage[],
    default: () => [],
  },
  participants: {
    type: Array as () => Participant[],
    default: () => [],
  },
});

defineEmits(["close-sidebar"]);
</script>
