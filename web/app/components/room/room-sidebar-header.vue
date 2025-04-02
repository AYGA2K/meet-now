<template>
  <div class="p-4">
    <div class="flex justify-between items-center">
      <h2 class="font-medium">
        {{ isChatOpen ? "Meeting chat" : "Participants (5)" }}
      </h2>
      <UButton
        icon="i-lucide-x"
        color="gray"
        variant="ghost"
        size="sm"
        aria-label="Close sidebar"
        @click="$emit('close-sidebar')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  isChatOpen: Boolean,
});

defineEmits(["close-sidebar"]);
</script>
