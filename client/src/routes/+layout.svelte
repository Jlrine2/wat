<script lang="ts">
	import '../app.css';
	let { children } = $props();
	import { authConfig, checkAuth } from '$lib/auth';
	
	import { Navbar, NavBrand, NavLi, NavUl, Avatar, Dropdown, DropdownHeader, DropdownItem, Button, Alert } from 'flowbite-svelte';
    import LoginWithDiscord from '$lib/components/loginWithDiscord.svelte';
</script>

{#await checkAuth()}
	<p>loading...</p>
{:then _} 
<Navbar>
	<NavBrand href="/">
			<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">WatchAnythingTogether</span>
	</NavBrand>
	{#if authConfig.isAuthenticated}
	  <div class="flex items-center md:order-2">
		<Avatar id="avatar-menu" src={authConfig.avatar} />
	  </div>
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
	  <div class="flex md:order-2">
		<LoginWithDiscord/>
	  </div>
	{/if}
	<NavUl>
	  <NavLi href="/">Home</NavLi>
	  <NavLi href="/watch">Watch</NavLi>
	</NavUl>
  </Navbar>

{@render children()}
{:catch}
<Alert color="red">
	<span class="font-medium">Something has gone Wrong!</span>
</Alert>
{/await}

