<template>
  <div class="p-4 border-t border-gray-200">
    <UInputGroup>
      <UInput
        v-model="newMessage"
        placeholder="Send a message to everyone"
        class="w-full"
      />
      <UButton
        icon="i-lucide-send"
        color="primary"
        variant="solid"
        :disabled="!newMessage"
        aria-label="Send message"
      />
    </UInputGroup>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const newMessage = ref("");
</script>
