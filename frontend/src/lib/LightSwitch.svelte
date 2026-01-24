<script lang="ts">
  import { Switch } from "@skeletonlabs/skeleton-svelte";
  import { onMount } from "svelte";

  let checked = false;

  // On mount, set checked state from localStorage (default: scheme-light)
  onMount(() => {
    const mode = localStorage.getItem("mode") || "scheme-light";
    checked = mode === "scheme-dark";
  });

  function onCheckedChange(event: { checked: boolean }) {
    const mode = event.checked ? "scheme-dark" : "scheme-light";
    document.documentElement.setAttribute("class", mode);
    localStorage.setItem("mode", mode);
    checked = event.checked;
  }
</script>

<svelte:head>
  <!-- Set mode on HTML tag ASAP to avoid FOUC -->
  <script>
    document.documentElement.setAttribute(
      "data-mode",
      localStorage.getItem("mode") || "scheme-light",
    );
  </script>
</svelte:head>

<Switch {checked} {onCheckedChange} aria-label="Toggle scheme-dark mode">
  <Switch.Control>
    <Switch.Thumb />
  </Switch.Control>
  <Switch.HiddenInput />
</Switch>
