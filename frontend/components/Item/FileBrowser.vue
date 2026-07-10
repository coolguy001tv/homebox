<template>
  <BaseModal v-model="open">
    <template #title> Browse NAS Files </template>

    <div class="flex flex-col gap-3 min-h-[300px] max-h-[60vh]">
      <!-- Breadcrumb navigation -->
      <div class="flex items-center gap-1 text-sm breadcrumbs overflow-x-auto">
        <ul class="flex items-center gap-0">
          <li>
            <button class="link link-hover" @click="goToDir(-1)">Home</button>
          </li>
          <li v-for="(crumb, idx) in breadcrumbs" :key="idx" class="flex items-center gap-0">
            <span class="mx-1">/</span>
            <button class="link link-hover truncate max-w-[150px]" @click="goToDir(idx)">
              {{ crumb }}
            </button>
          </li>
        </ul>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center items-center flex-1">
        <span class="loading loading-spinner loading-md"></span>
      </div>

      <!-- Error -->
      <div v-else-if="errorMsg" class="alert alert-error flex-1">
        <Icon name="mdi-alert-circle" class="h-5 w-5" />
        <span>{{ errorMsg }}</span>
      </div>

      <!-- Empty -->
      <div v-else-if="entries.length === 0" class="flex justify-center items-center flex-1 text-gray-500">
        <div class="text-center">
          <Icon name="mdi-folder-open-outline" class="h-12 w-12 mx-auto mb-2" />
          <p>This directory is empty</p>
        </div>
      </div>

      <!-- File/directory grid -->
      <div v-else class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-3 overflow-y-auto flex-1 p-1">
        <button
          v-for="entry in sortedEntries"
          :key="entry.path"
          class="flex flex-col items-center gap-1 p-2 rounded-lg hover:bg-base-200 transition-colors text-center cursor-pointer border border-transparent hover:border-base-300"
          :class="{ 'opacity-60': !entry.isDir && !entry.isImage }"
          @click="onEntryClick(entry)"
        >
          <!-- Folder icon -->
          <Icon v-if="entry.isDir" name="mdi-folder" class="h-10 w-10 text-yellow-500" />

          <!-- Image thumbnail -->
          <img
            v-else-if="entry.isImage"
            :src="thumbSrc(entry.path)"
            class="h-16 w-full object-cover rounded"
            loading="lazy"
            :alt="entry.name"
          />

          <!-- Generic file icon -->
          <Icon v-else name="mdi-file-outline" class="h-10 w-10 text-gray-400" />

          <span class="text-xs break-all line-clamp-2 leading-tight">{{ entry.name }}</span>
          <span v-if="!entry.isDir" class="text-xs text-gray-400">{{ formatSize(entry.size) }}</span>
        </button>
      </div>
    </div>

    <template v-if="!loading && !errorMsg" #footer>
      <div class="text-xs text-gray-400">
        {{ entries.length }} item{{ entries.length !== 1 ? 's' : '' }}
      </div>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import type { FileEntry } from "~~/lib/api/classes/import";

type Props = {
  modelValue: boolean;
};

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
});

const emit = defineEmits<{
  (e: "update:modelValue", v: boolean): void;
  (e: "select", entry: FileEntry): void;
}>();

const open = useVModel(props, "modelValue", emit);

const api = useUserApi();
const toast = useNotifier();

const loading = ref(false);
const errorMsg = ref("");
const entries = ref<FileEntry[]>([]);
const currentPath = ref("");

// Breadcrumbs: split path segments
const breadcrumbs = computed(() => {
  if (!currentPath.value) return [];
  return currentPath.value.split("/").filter(Boolean);
});

const sortedEntries = computed(() => {
  return [...entries.value].sort((a, b) => {
    // Directories first
    if (a.isDir !== b.isDir) return a.isDir ? -1 : 1;
    // Alphabetical within each group
    return a.name.localeCompare(b.name);
  });
});

function thumbSrc(filePath: string): string {
  return api.items.importDir.thumbURL(filePath, 120);
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

async function loadDir(subPath: string) {
  loading.value = true;
  errorMsg.value = "";

  const { data, error } = await api.items.importDir.browse(subPath);

  loading.value = false;

  if (error) {
    errorMsg.value = `Failed to browse directory: ${error}`;
    entries.value = [];
    return;
  }

  entries.value = data || [];
}

function goToDir(index: number) {
  if (index < 0) {
    currentPath.value = "";
  } else {
    const crumbs = breadcrumbs.value.slice(0, index + 1);
    currentPath.value = crumbs.join("/");
  }
  loadDir(currentPath.value);
}

function onEntryClick(entry: FileEntry) {
  if (entry.isDir) {
    // Navigate into directory — strip the base prefix to get a relative subPath
    const basePrefix = getBasePrefix();
    if (entry.path.startsWith(basePrefix)) {
      currentPath.value = entry.path.substring(basePrefix.length);
    }
    loadDir(currentPath.value);
  } else {
    // Import the file
    emit("select", entry);
  }
}

function getBasePrefix(): string {
  // The browse API returns absolute paths. We need to find the base prefix
  // from the first entry so we can construct relative paths for navigation.
  if (entries.value.length === 0) return "";
  const firstPath = entries.value[0].path;
  const lastName = entries.value[0].name;
  return firstPath.substring(0, firstPath.length - lastName.length);
}

// Reload on open
watch(open, async (isOpen) => {
  if (isOpen) {
    currentPath.value = "";
    await loadDir("");
  }
});
</script>
