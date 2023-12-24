import { z } from "zod";


const post_schema = z.object({
    id: z.number(),
    slug: z.string(),
    title: z.string(),
    content: z.string(),
    category: z.string(),
    archived: z.union([z.literal(0), z.literal(1)]),
    publish_at: z.string(),
    created_at: z.string(),
    updated_at: z.string()
});

export type Post = z.infer<typeof post_schema>;


export const SCHEMA = {
    POST: post_schema
}


