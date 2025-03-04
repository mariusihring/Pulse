'use client'
import { useState, useEffect } from "react";
import { useMutation, useSubscription } from "urql";
import { graphql } from "@/lib/gql";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {Progress} from "@/components/ui/progress";

export default function Home() {
    return (
        <p>Index</p>

    );
}