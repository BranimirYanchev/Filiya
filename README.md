# Filia Project Frontend

This is the frontend part of the **Filia Project** — a React application styled with SCSS and structured with modular components.

## 🚀 Getting Started

### Prerequisites
- Node.js (16+ recommended)
- npm or yarn

### Installation
Clone the repository and switch to the project folder:
```bash
git clone https://github.com/Marionvd/filia-project-frontend.git
cd filia-project-frontend
git checkout home-page
```

Install dependencies:
```bash
npm install
# or
yarn
```

### Running the Project
Start the development server:
```bash
npm start
# or
yarn start
```
Once started, the app will be available at [http://localhost:3000](http://localhost:3000).

### Building for Production
```bash
npm run build
# or
yarn build
```
This creates an optimized build inside the `build/` directory, ready for deployment.

## 📂 Project Structure
```
src/
├── components/         # Reusable components
├── pages/              # Page-level views
├── styles/             # SCSS modules and global styles
├── assets/             # Images, fonts, and static assets
├── App.jsx             # Root component
└── index.js            # Entry point
public/
├── index.html
└── ...
package.json
README.md
```

## 📜 Available Scripts
| Command          | Description                                |
|------------------|--------------------------------------------|
| `npm start`      | Run the development server with hot reload |
| `npm run build`  | Create a production build                  |
| `npm test`       | Run tests (if available)                   |
| `npm run lint`   | Run linter/formatter (if configured)       |

## 🌐 Deployment
You can deploy this project using:
- **Netlify** → Build command: `npm run build`, Publish directory: `build/`
- **Vercel** → Framework: React, Output directory: `build/`
- **GitHub Pages** → Run `npm run build` and deploy the contents of the `build/` folder

---

💡 With this setup you can quickly run, develop, and deploy the Filia Project frontend.
