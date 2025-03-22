import type { PageLoad } from './$types';
export interface AuthConfig {
    avatar: string 
    username: string | null;
    isAuthenticated: boolean;
}

export let authConfig: AuthConfig = {
    avatar: "",
    username: null,
    isAuthenticated: false

};

export async function checkAuth() {
    try {
        const response = await fetch(`${document.location.origin}/auth/me`);
        console.log("auth resp: ", response.status)
        if (response.ok) {
            const data = await response.json();
            let avatar_suffix = ".png"
            if ((data.user.avatar as string).startsWith('a_')) {
                avatar_suffix = ".gif"
            }
            authConfig = {
                avatar: `https://cdn.discordapp.com/avatars/${data.user.id}/${data.user.avatar}.png`,
                username: data.user.global_name,
                isAuthenticated: true
            };
        } else {
            authConfig = {
                avatar: "",
                username: null,
                isAuthenticated: false
            };
        }
    } catch (error) {
        console.error('Auth check failed:', error);
        authConfig = {
            avatar: "",
            username: null,
            isAuthenticated: false
        };
    }
}