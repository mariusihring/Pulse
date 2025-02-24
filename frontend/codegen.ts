import type { CodegenConfig } from '@graphql-codegen/cli'
 
const config: CodegenConfig = {
   schema: 'http://localhost:3001/query',
   documents: ['src/**/*.vue', "pages/**/*.vue"],
   generates: {
      './src/gql/': {
        preset: 'client',
        config: {
            useTypeImports: true
        }
      }
   }
}
export default config