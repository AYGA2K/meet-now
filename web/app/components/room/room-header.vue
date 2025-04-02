<template>
  <header class="shadow-sm p-2">
    <div class="container mx-auto flex justify-between items-center">
      <div class="flex items-center space-x-2">
        <UIcon name="i-lucide-video" class="text-emerald-600 text-xl" />
        <h1 class="text-xl font-semibold">MeetNow</h1>
      </div>
      <div class="flex items-center space-x-4">
        <UButton
          icon="i-lucide-layout-grid"
          color="gray"
          variant="ghost"
          aria-label="Change layout"
        />
        <UButton
          icon="i-lucide-settings"
          color="gray"
          variant="ghost"
          aria-label="Settings"
        />
        <UAvatar src="https://i.pravatar.cc/150?img=68" alt="User" size="sm" />
      </div>
    </div>
  </header>
</template>
