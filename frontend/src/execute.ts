import axios, { AxiosError } from "axios";
import { supabase } from "./lib/supabase";
import { TypedDocumentString } from "./graphql/graphql";

interface GraphQLError {
	message: string;
	locations: Array<{ line: number; column: number }>;
	path: string[];
}

interface GraphQLResponse<T> {
	data?: T;
	errors?: GraphQLError[];
}

class GraphQLErrorsError extends Error {
	constructor(public errors: GraphQLError[]) {
		super(errors.map((e) => e.message).join(", "));
		this.name = "GraphQLErrorsError";
	}
}

export async function execute<TResult, TVariables>(
  query: TypedDocumentString<TResult, TVariables>,
  variables?: TVariables
): Promise<TResult> {
  try {
    const { data: { session } } = await supabase.auth.getSession();
    const token = session?.access_token;

    const response = await axios<GraphQLResponse<TResult>>({
      url: "http://127.0.0.1:8080/query",
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/graphql-response+json",
        Authorization: token ? `Bearer ${token}` : "",
      },
      data: {
        query,
        variables,
      },
    });

    if (response.data.errors) {
      throw new GraphQLErrorsError(response.data.errors);
    }

    return response.data.data!;
  }  catch (error) {
	if (error instanceof GraphQLErrorsError) {
		throw error;
	}
	if (axios.isAxiosError(error)) {
		const axiosError = error as AxiosError<GraphQLResponse<TResult>>;
		if (axiosError.response?.data.errors) {
			throw new GraphQLErrorsError(axiosError.response.data.errors);
		}
	}
	throw error;
}
}
