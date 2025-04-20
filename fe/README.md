# Changi Airport Weather Report System - Frontend

This is the frontend application for the Changi Airport Weather Report System, a full-stack application for generating, storing, and comparing weather reports for Changi Airport using the OpenWeatherMap API.

## Technology Stack

- **Framework**: [Next.js 15](https://nextjs.org) with App Router
- **UI Components**: [Shadcn UI](https://ui.shadcn.com/) - A collection of reusable components built with Radix UI and Tailwind CSS
- **Styling**: [Tailwind CSS](https://tailwindcss.com/) with custom configuration
- **Form Handling**: [React Hook Form](https://react-hook-form.com/) with [Zod](https://zod.dev/) validation
- **Notifications**: [Sonner](https://sonner.emilkowal.ski/) for toast notifications
- **Icons**: [Lucide React](https://lucide.dev/) for beautiful icons
- **Fonts**: [Geist](https://vercel.com/font) font family from Vercel
- **Testing**: [Jest](https://jestjs.io/) and [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)

## Features

- **Weather Report Generation**: Generate weather reports for Changi Airport for the current time or a specific date
- **History View**: Browse through all previously generated reports with pagination
- **Report Comparison**: Select two reports to compare their weather data and see deviations
- **Responsive Design**: Fully responsive UI that works on mobile and desktop
- **Form Validation**: Client-side form validation for all inputs
- **Error Handling**: Proper error handling and user feedback

## Project Structure

```
fe/
├── public/              # Static assets
├── src/
│   ├── app/             # Next.js App Router pages
│   │   ├── compare/     # Report comparison page
│   │   ├── history/     # Report history page
│   │   └── page.tsx     # Home page (report generation)
│   ├── components/      # Reusable UI components
│   │   └── ui/          # Shadcn UI components
│   └── lib/             # Utility functions and API clients
│       ├── api/         # API client code
│       │   ├── client/  # API client functions
│       │   ├── constants/ # API constants
│       │   └── types/   # TypeScript types for API
│       └── utils.ts     # Utility functions
├── next.config.js       # Next.js configuration
├── package.json         # Dependencies and scripts
└── tsconfig.json        # TypeScript configuration
```

## Getting Started

### Prerequisites

- Node.js 18 or later
- npm or yarn

### Installation

1. Install dependencies:

```bash
npm install
# or
yarn install
```

2. Set up environment variables:

Create a `.env.local` file in the root directory with the following content:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### Development

Run the development server:

```bash
npm run dev
# or
yarn dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the application.

### Testing

Run the tests:

```bash
npm test
# or
yarn test
```

Run the tests in watch mode:

```bash
npm run test:watch
# or
yarn test:watch
```

Generate test coverage report:

```bash
npm run test:coverage
# or
yarn test:coverage --ignore-engines
```

### Building for Production

```bash
npm run build
npm run start
# or
yarn build
yarn start
```

## Docker Integration

This frontend is designed to work with Docker and is configured in the project's `docker-compose.yml` file. The frontend container:

- Builds from the `./fe` directory using the provided Dockerfile
- Runs on port 3000
- Connects to the backend service
- Uses the environment variable `NEXT_PUBLIC_API_URL` to communicate with the API

To run the entire application (frontend, backend, and MongoDB) using Docker:

```bash
docker-compose up -d
```

## Learn More

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API
- [Shadcn UI Documentation](https://ui.shadcn.com/docs) - learn about the UI components
- [Tailwind CSS Documentation](https://tailwindcss.com/docs) - learn about Tailwind CSS
- [Jest Documentation](https://jestjs.io/docs/getting-started) - learn about Jest testing framework
- [React Testing Library Documentation](https://testing-library.com/docs/react-testing-library/intro/) - learn about React Testing Library
