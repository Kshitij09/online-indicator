// Define the User interface to match the API response
export interface User {
    id: string;
    name: string; // We'll map username to name for backward compatibility
    online: boolean; // We'll map is_online to online for backward compatibility
    last_online?: number; // Optional timestamp
}

// Define the API response interface
interface UserStatus {
    id: string;
    name: string;
    is_online: boolean;
    last_online: number;
}

interface StatusResponse {
    items: UserStatus[];
}

// List of user IDs to fetch
const userIds = Array.from({length: 100}, (_, i) => (i + 1).toString());

/**
 * Fetch users from the API
 * @returns Promise that resolves with the users
 */
export const fetchUsers = async (): Promise<User[]> => {
    try {
        // Make API call to fetch users via proxy to avoid CORS issues
        const response = await fetch('/api/batch/status', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                ids: userIds
            }),
        });

        if (!response.ok) {
            throw new Error(`API call failed with status: ${response.status}`);
        }

        // Parse the response
        const data: StatusResponse = await response.json();

        // Map API response to User interface
        return data.items.map(item => ({
            id: item.id,
            name: item.name, // Map username to name
            online: item.is_online, // Map is_online to online
            last_online: item.last_online
        }));
    } catch (error) {
        console.error('Error fetching users:', error);
        // Return empty array in case of error
        return [];
    }
};
