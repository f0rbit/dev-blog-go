import { z } from "zod";

const post_schema = z.object({
    id: z.number(),
    slug: z.string(),
    title: z.string(),
    content: z.string(),
    category: z.string(),
    author_id: z.number(),
    tags: z.array(z.string()),
    archived: z.boolean(),
    publish_at: z.string(),
    created_at: z.string(),
    updated_at: z.string()
});


const posts_response_schema = z.object({
    posts: z.array(post_schema),
    total_posts: z.number(),
    total_pages: z.number(),
    per_page: z.number(),
    current_page: z.number(),
})

const category_schema = z.object({
    name: z.string(),
    parent: z.string()
});

const category_node_schema = z.object({
    name: z.string(),
    children: z.array(z.any())
})

const category_response = z.object({
    graph: category_node_schema,
    categories: z.array(category_schema)
});

const access_key = z.object({
    id: z.number(),
    value: z.string(),
    user_id: z.number(),
    name: z.string(),
    note: z.string(),
    enabled: z.boolean(),
    created_at: z.string(),
    updated_at: z.string(),
});

const integration_link = z.object({
    id: z.number(),
    user_id: z.number(),
    last_fetch: z.string(),
    location: z.string(),
    source: z.string(),
    data: z.any(),
    created_at: z.string(),
    updated_at: z.string(),
});

export type Post = z.infer<typeof post_schema>;

export type PostsResponse = z.infer<typeof posts_response_schema>;

export type AccessKey = z.infer<typeof access_key>;

export type Category = z.infer<typeof category_schema>;

export type IntegrationLink = z.infer<typeof integration_link>;

// i don't think there's a way for zod to do self-recursive types yet.
export type CategoryNode = { name: string, children: CategoryNode[] };

export type CategoryResponse = {
    categories: Category[]
    graph: CategoryNode
};


export const SCHEMA = {
    POST: post_schema,
    POSTS_RESPONSE: posts_response_schema,
    CATEGORY: category_schema,
    CATEGORY_NODE: category_node_schema,
    CATEGORY_RESPONSE: category_response,
    ACCESS_KEY: access_key
}


