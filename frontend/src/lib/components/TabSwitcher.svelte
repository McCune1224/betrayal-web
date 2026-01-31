<script lang="ts">
  interface Tab {
    id: string;
    label: string;
    activeColor: string;
  }

  interface Props {
    activeTab: string;
    tabs: Tab[];
    onTabChange: (tabId: string) => void;
  }

  let { activeTab, tabs, onTabChange }: Props = $props();
</script>

<div class="flex border-b border-surface-200-800 relative">
  {#each tabs as tab}
    <button
      class="flex-1 py-4 text-center font-bold tracking-wide text-sm uppercase transition-all duration-200 {activeTab === tab.id
        ? `bg-surface-100-900 text-${tab.activeColor}-600-400 border-b-2 border-${tab.activeColor}-500`
        : 'bg-surface-200-800 text-surface-500 hover:text-surface-700-300 hover:bg-surface-300-700'}"
      onclick={() => onTabChange(tab.id)}
    >
      {tab.label}
    </button>
  {/each}
</div>
