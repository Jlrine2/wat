<script lang="ts">
	import '../app.css';
	let { children } = $props();
	import { authConfig, checkAuth } from '$lib/auth';
	
	import { Navbar, NavBrand, Avatar, Dropdown, DropdownHeader, DropdownItem, Frame } from 'flowbite-svelte';
    import LoginWithDiscord from '$lib/components/loginWithDiscord.svelte';
</script>


<Navbar color="gray" navContainerClass="flex flex-wrap justify-between -container items-center mx-0 w-full">
	<NavBrand href="/">
		<img src="/logo.svg" class="me-3 h-6 sm:h-9" alt="Watch Anything Together Logo" />
		<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">WatchAnythingTogether</span>
	</NavBrand>

	{#await checkAuth()}
	<Avatar />
	{:then _} 
	{#if authConfig.isAuthenticated}
	<Avatar id="avatar-menu" src={authConfig.avatar} />
	<Dropdown placement="bottom" triggeredBy="#avatar-menu">
	<DropdownHeader>
		<span class="block text-sm">{authConfig.username}</span>
	</DropdownHeader>
	<DropdownItem on:click={() => {
		document.cookie = "watAuth=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
		window.location.reload();
	}}>Sign out</DropdownItem>
	</Dropdown>
	{:else}
	<LoginWithDiscord/>
	{/if}
	{/await}
  </Navbar>

{@render children()}

