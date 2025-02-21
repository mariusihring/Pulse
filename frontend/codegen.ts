import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: "http://0.0.0.0:3001/query",
  documents: ["src/**/*.{ts,tsx}"],
  ignoreNoDocuments: true,
  generates: {
    "./src/graphql/": {
      preset: "client",
      plugins: ["typescript", "typescript-operations"],
      config: {
        documentMode: "string",
        dedupeFragments: true,
        skipTypename: false,
        enumsAsTypes: true,
        // Generate types for all schema types, not just the ones used in operations
        onlyOperationTypes: false,
      },
    },
    "./schema.graphql": {
      plugins: ["schema-ast"],
      config: {
        includeDirectives: true,
      },
    },
    // Add a new output specifically for base types
    // "./src/graphql/types.ts": {
    //   plugins: ["typescript"],
    //   config: {
    //     // This ensures we get types for everything in the schema
    //     enumsAsTypes: true,
    //     skipTypename: false,
    //     // Include scalars configuration if you have custom scalars
    //     scalars: {
    //       DateTime: "string",
    //       // Add other custom scalars here
    //     },
    //   },
    // },
  },
};

export default config;
