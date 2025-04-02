<template>
  <div class="flex-1 overflow-y-auto p-4 space-y-4">
    <div
      v-for="(message, index) in chatMessages"
      :key="index"
      class="flex items-start"
    >
      <UAvatar
        :src="`https://i.pravatar.cc/150?img=${message.avatar}`"
        :alt="message.sender"
        size="sm"
        class="mr-2 flex-shrink-0"
      />
      <div>
        <div class="flex items-center">
          <span class="font-medium text-sm">{{ message.sender }}</span>
          <span class="ml-2 text-xs">{{ message.time }}</span>
        </div>
        <p class="text-sm">{{ message.text }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import type { ChatMessage } from "~/types";

defineProps({
  chatMessages: {
    type: Array as () => ChatMessage[],
    required: true,
  },
});
</script>
