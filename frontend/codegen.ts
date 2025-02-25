import type {CodegenConfig} from "@graphql-codegen/cli"

const config: CodegenConfig = {
    schema: "http://localhost:3001/query",
    documents: ["./app/**/*.tsx"],
    generates: {
        "./lib/gql/": {
            preset: "client"
        }
    }
}

export default config