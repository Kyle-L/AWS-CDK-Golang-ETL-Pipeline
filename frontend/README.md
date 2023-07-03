This is a [Next.js](https://nextjs.org/) project bootstrapped with [`create-next-app`](https://github.com/vercel/next.js/tree/canary/packages/create-next-app). It is the frontend that is used to visualize banking data loaded into DynamoDB and accessible through API Gateway.

## Getting Started

First, create a new `.env.local` file with the following format in the root of the `/frontend` directory.

```
API_ENDPOINT=https://<your-api-gw-endpoint>
```

Then, run the development server:

```bash
npm run dev
# or
yarn dev
```

Now, open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `pages/index.tsx`. The page auto-updates as you edit the file.

## Deployment

Run the following command to build the production ready application.

```bash
npm run build
#
yarn build
```

At this point, you can upload all files to the frontend S3 bucket that was deployed as part of the CDK project.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!