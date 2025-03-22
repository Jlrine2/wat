import { checkAuth } from "$lib/auth"
import { redirect } from "@sveltejs/kit"

export const ssr = false
export const prerender = true;
export const trailingSlash = 'always';