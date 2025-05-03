<script lang="ts">
  import {onDestroy, onMount} from 'svelte';
  import {writable} from 'svelte/store';
  import {fetchUsers, type User} from './services/UserService';

  // Create a writable store for users in the component
    const usersStore = writable<User[]>([]);

    // Search functionality
    let searchTerm = '';

    // Derived filtered users based on search term
    $: filteredUsers = searchTerm
        ? $usersStore.filter(user =>
            user.name.toLowerCase().includes(searchTerm.toLowerCase()))
        : $usersStore;

    // Calculate stats
    $: totalUsers = $usersStore.length;
    $: onlineUsers = $usersStore.filter(user => user.online).length;

    // Loading state
    let loading = true;

    // Variable to store the interval ID
    let pollingInterval: number;

    // Function to fetch users and update the store
    async function fetchAndUpdateUsers() {
        try {
            // Fetch users from the service
            const users = await fetchUsers();
            // Update the store with fetched users
            usersStore.set(users);
        } catch (error) {
            console.error('Error fetching users:', error);
        } finally {
            if (loading) {
                loading = false;
            }
        }
    }

    // Initialize the component
    onMount(async () => {
        // Initial fetch
        await fetchAndUpdateUsers();

        // Set up polling every 5 seconds (5000 ms)
        pollingInterval = setInterval(fetchAndUpdateUsers, 5000);
    });

    // Clean up when component is destroyed
    onDestroy(() => {
        if (pollingInterval) {
            clearInterval(pollingInterval);
        }
    });
</script>

<div class="users-container">
    <h2>Users</h2>

    {#if loading}
        <div class="loading-container">
            <div class="loading-spinner"></div>
            <p>Loading users...</p>
        </div>
    {:else if $usersStore.length === 0}
        <p>No users available</p>
    {:else}
        <div class="users-stats">
            <span>Total: {totalUsers}</span>
            <span>Online: {onlineUsers}</span>
            <span>Offline: {totalUsers - onlineUsers}</span>
        </div>

        <div class="search-container">
            <input
                    type="text"
                    placeholder="Search users..."
                    bind:value={searchTerm}
                    class="search-input"
            />
        </div>

        <div class="users-list-container">
            {#if filteredUsers.length === 0}
                <p class="no-results">No users match your search criteria</p>
            {:else}
                <ul class="users-list">
                    {#each filteredUsers as user (user.id)}
                        <li class="user-item">
              <span
                      class="status-indicator"
                      class:online={user.online}
                      class:offline={!user.online}
              ></span>
                            <span class="user-name">{user.name}</span>
                        </li>
                    {/each}
                </ul>
            {/if}
        </div>
    {/if}
</div>

<style>
    /* Define CSS variables for light/dark mode */
    :root {
        --stats-bg-color: #f9f9f9;
        --stats-text-color: #1a1a1a;
        --no-results-color: #888;
        --spinner-border-color: rgba(0, 0, 0, 0.3);
    }

    @media (prefers-color-scheme: dark) {
        :root {
            --stats-bg-color: #1a1a1a;
            --stats-text-color: rgba(255, 255, 255, 0.87);
            --no-results-color: #aaa;
            --spinner-border-color: rgba(255, 255, 255, 0.3);
        }
    }

    /* Override global centering */
    :global(body) {
        display: block !important;
        place-items: unset !important;
    }

    :global(#app) {
        max-width: 100% !important;
        padding-top: 1rem !important;
    }

    .users-container {
        width: 100%;
        max-width: 800px;
        margin: 0 auto;
        text-align: left;
        padding: 20px;
    }

    /* Stats section */
    .users-stats {
        display: flex;
        justify-content: space-between;
        margin-bottom: 15px;
        padding: 10px;
        border-radius: 4px;
        background-color: var(--stats-bg-color, #f9f9f9);
        color: var(--stats-text-color, #1a1a1a);
    }

    .users-stats span {
        font-size: 14px;
    }

    /* Search section */
    .search-container {
        margin-bottom: 15px;
        padding: 10px;
        border-radius: 4px;
        background-color: var(--stats-bg-color, #f9f9f9);
    }

    .search-input {
        width: 100%;
        padding: 8px 12px;
        border: none;
        background-color: transparent;
        color: var(--stats-text-color, #1a1a1a);
        box-sizing: border-box;
    }

    .search-input:focus {
        outline: none;
    }

    /* Users list */
    .users-list-container {
        max-height: 70vh;
        overflow-y: auto;
        border: 1px solid #999;
        border-radius: 4px;
    }

    .users-list {
        list-style: none;
        padding: 0;
        margin: 0;
    }

    .user-item {
        display: flex;
        align-items: center;
        padding: 10px;
        border-bottom: 1px solid #999;
    }

    .user-item:last-child {
        border-bottom: none;
    }

    .status-indicator {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        margin-right: 10px;
        cursor: pointer;
    }

    .online {
        background-color: #4CAF50; /* Green */
    }

    .offline {
        background-color: #F44336; /* Red */
    }

    .user-name {
        font-size: 16px;
    }

    .no-results {
        text-align: center;
        padding: 20px;
        color: var(--no-results-color, #888);
    }

    .loading-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 40px 0;
    }

    .loading-spinner {
        width: 30px;
        height: 30px;
        border: 3px solid rgba(255, 255, 255, 0.3);
        border-radius: 50%;
        border-top-color: #646cff;
        border-bottom-color: #646cff;
        animation: spin 1s ease-in-out infinite;
        margin-bottom: 10px;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    /* Responsive adjustments */
    @media (max-width: 600px) {
        .users-stats {
            flex-direction: column;
            gap: 5px;
        }

        .users-container {
            padding: 10px;
        }
    }
</style>
