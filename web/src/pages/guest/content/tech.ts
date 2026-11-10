export interface TechCategory {
  category: string
  description: string
  items: {
    name: string
    version?: string
    description: string
    badges: string[]
  }[]
}

export interface TechContent {
  headline: string
  title: string
  summary: string
  categories: TechCategory[]
}

export const techContent: TechContent = {
  headline: 'Project Architecture & Stack',
  title: 'Technology Stack',
  summary:
    'A robust full-stack web application combining a high-performance Go backend with a modern Vue 3 frontend, styled with a distinctive neo-brutalist design system.',
  categories: [
    {
      category: 'Backend & API',
      description: 'High-performance Go services providing RESTful endpoints, database migrations, and authentication.',
      items: [
        {
          name: 'Go',
          version: '1.27.3',
          description: 'Core backend language offering exceptional concurrency, fast execution, and strong typing.',
          badges: ['Language', 'Performance'],
        },
        {
          name: 'Fiber v3',
          version: 'v3.6.1',
          description:
            'Express-inspired web framework built on top of fasthttp for lightning-fast HTTP request routing.',
          badges: ['Web Framework', 'REST API'],
        },
        {
          name: 'GORM & pgx',
          version: 'v1.31.2 / v5',
          description: 'Robust Object-Relational Mapping with the high-performance PostgreSQL driver for Go.',
          badges: ['ORM', 'Database Driver'],
        },
        {
          name: 'golang-migrate',
          version: 'v4.20.1',
          description: 'Database migration management tool ensuring version-controlled schema evolution.',
          badges: ['Migrations'],
        },
        {
          name: 'JWT Authentication',
          version: 'v5.3.1',
          description: 'Secure token-based authentication for managing user sessions and protected API routes.',
          badges: ['Auth', 'Security'],
        },
        {
          name: 'AWS SDK for Go v2',
          version: 'v1.47.8',
          description: 'S3-compatible client SDK used to interface with Garage object storage for file uploads.',
          badges: ['Storage', 'S3 Compatible'],
        },
        {
          name: 'Cobra & Godotenv',
          version: 'v1.10.2',
          description: 'CLI command management and environment configuration loader.',
          badges: ['CLI', 'Config'],
        },
      ],
    },
    {
      category: 'Frontend & UI',
      description: 'Reactive client-side SPA built with Vue 3, TypeScript, Pinia, and Tailwind CSS.',
      items: [
        {
          name: 'Vue 3',
          version: '^3.5.45',
          description: 'Progressive JavaScript framework utilizing Composition API for reactive UI components.',
          badges: ['Frontend', 'Composition API'],
        },
        {
          name: 'Vite',
          version: '^8.3.6',
          description: 'Lightning-fast frontend build tool and development server with hot module replacement.',
          badges: ['Build Tool', 'Bundler'],
        },
        {
          name: 'TypeScript',
          version: '^6.0.3',
          description: 'Strict type safety across components, stores, and API clients.',
          badges: ['Language', 'Type Safety'],
        },
        {
          name: 'Pinia',
          version: '^4.1.0',
          description: 'Intuitive, type-safe state management store for Vue applications.',
          badges: ['State Management'],
        },
        {
          name: 'Vue Router',
          version: '^5.4.0',
          description: 'Official router for Vue.js handling client-side page routing and layouts.',
          badges: ['Routing'],
        },
        {
          name: 'Tailwind CSS v4',
          version: '^4.3.5',
          description: 'Utility-first CSS framework styled with custom neo-brutalist themes and brutal components.',
          badges: ['Styling', 'Neo-Brutalist'],
        },
        {
          name: 'Vitest & Vue Test Utils',
          version: '^5.1.0',
          description: 'Fast unit testing framework built on Vite for testing Vue components and logic.',
          badges: ['Testing'],
        },
      ],
    },
    {
      category: 'Database & Infrastructure',
      description: 'Persistent data storage and tooling for local and production environments.',
      items: [
        {
          name: 'PostgreSQL',
          description: 'Powerful, open-source relational database managing structured application data.',
          badges: ['Database', 'SQL'],
        },
        {
          name: 'Garage Storage',
          description: 'Distributed S3-compatible object storage used for reliable file storage and assets.',
          badges: ['Object Storage', 'Garage S3'],
        },
      ],
    },
    {
      category: 'Developer Tooling & AI',
      description: 'Advanced development intelligence, MCP servers, and automated workflows.',
      items: [
        {
          name: 'CodeGraph',
          description: 'SQLite knowledge graph indexing workspace symbols, calls, and code structure.',
          badges: ['Code Intelligence', 'AST'],
        },
        {
          name: 'Model Context Protocol (MCP)',
          description: 'Standardized protocol enabling AI agents to interact with application databases and tools.',
          badges: ['MCP', 'AI Integration'],
        },
      ],
    },
  ],
}
