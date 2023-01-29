# appsyncgen [![Go Report Card](https://goreportcard.com/badge/github.com/kopkunka55/appsyncgen)](https://goreportcard.com/report/github.com/kopkunka55/appsyncgen) [![Go Reference](https://pkg.go.dev/badge/github.com/kopkunka55/appsyncgen.svg)](https://pkg.go.dev/github.com/kopkunka55/appsyncgen)

appsyncgen is a CLI for generating AWS AppSync JavaScript Resolvers based on Amazon DynamoDB single-table design.

appsyncgen is inspired by [AWS Amplify CLI](https://docs.amplify.aws/cli/)

## Overview

appsyncgen is a CLI providing some useful capability to develop GraphQL API with AppSync JavaScript Resolver using Amazon DynamoDB single table.

**appsyncgen** provides:

* Generate JS resolvers from `schema.graphql`.
* Support some directive supported by AWS Amplify (`@auth`, `@hasOne`, `@hasMany`, `@manyToMany`).
* `@auth` directive supports multiple providers (`apiKey`, `oidc`, `iam`, `lambda`, `userPools`) which provide resolver-level authorization.
* Export resolver list by JSON so that you can easily implement CDK stack.
* Generate CloudFormation Template automatically.
* Generate Pipeline resolvers for your queries and mutations so that you can slot in your custom business logic between generated resolvers.
* Optimistic locking with version number for update resolver
* **TypeScript** support is coming soon.

## Concepts

[AWS Amplify CLI](https://docs.amplify.aws/cli/) is really powerful tool to generate AppSync resolvers using multi-table Amazon DynamoDB, but sometimes building complicated logic is not easy with VTL. **appsyncgen** helps us generate **JavaScript** Resolver, which should make it easier to edit auto-generated code. Additionally **appsyncgen** generate resolvers based on Amazon DynamoDB [single-table design](https://aws.amazon.com/blogs/compute/creating-a-single-table-design-with-amazon-dynamodb/), so that you don't need to consider DynamoDB key design by yourself.

## Installing
```shell
go install github.com/kopkunka55/appsyncgen
```
## Usage
All you need is `schema.graphql` which includes only basic types. Mutation/Query/Subscription and some supplemental types will be added.

```graphql
type Message
@auth (rules: [
    {provider: apiKey},
])
{
    body: String!
    from: User! @hasOne
}

enum Role {
  ADMIN
  READER
  EDITOR
}

type User
@auth (rules: [
    {provider: apiKey, operations: [create, update, delete, read]},
])
{
    name: String!
    chats: [Chat]! @manyToMany
    profilePicture: String
    roles: Role!
}

type Chat
@auth (rules: [
    {provider: apiKey}
])
{
    name: String!
    members: [User] @manyToMany
    messages: [Message] @hasMany
}
```

```shell
appsyncgen generate --output='./resolvers' --schema='./schema.graphql' --name='appsyncgen'
```

## License

The source code for the site is licensed under the MIT license, which you can find in the [LICENSE](./LICENSE) file.


