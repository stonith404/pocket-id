import React from "react"; // Font Awesome// Ensure Tailwind is imported in your global styles
import "/styles.css";
import "@fortawesome/fontawesome-free/css/all.min.css";

const Home: React.FC = () => {
  return (
    <div className="text-white h-screen flex flex-col bg-muted/40" >
      <header className="" style={{ backgroundColor: "hsl(240, 10%, 3.9%)" }}>
        <div className="w-full border-b">
          <div className="max-w-[1640px] mx-auto flex w-full items-center justify-between px-4 md:px-10">
            <div className="flex h-16 items-center justify">
              <img src="https://docs.pocket-id.org/img/pocket-id.png" alt="Pocket ID Logo" className="mr-3 h-8 w-8" />
              <h1 className="text-lg font-medium" style={{ margin: 0 }} >Pocket ID</h1>
            </div>
            <a href="https://github.com/stonith404/pocket-id" target="_blank" rel="noopener noreferrer" style={{ color: "hsl(0, 0%, 98%)" }} className="text-white text-2xl">
              <i className="fab fa-github" aria-hidden="true"></i>
            </a>
          </div>
        </div>
      </header>

      <main className="flex-1 flex flex-col justify-center items-center px-4 sm:px-0">
        <section className="text-center py-12 sm:py-20 glass m-4 w-full max-w-4xl animate-fadeIn px-6 sm:px-12">
          <img src="https://docs.pocket-id.org/img/pocket-id.png" alt="Pocket ID Logo" className="h-24 w-24 mx-auto mb-6 animate-bounce" />
          <h2 className="text-3xl sm:text-4xl font-extrabold">Secure Your Services</h2>
          <p className="mt-4 text-lg">A simple, open-source OIDC provider leveraging passkeys for secure authentication.</p>
          <a href="https://docs.pocket-id.org" target="_blank" rel="noopener noreferrer" className="mt-6 inline-block text-black px-6 py-3 rounded-lg font-semibold" style={{ backgroundColor: "hsl(0, 0%, 98%)", color: "hsl(240, 10%, 3.9%)"}}>
            Get Started
          </a>
        </section>

        <section className="container mx-auto py-12 sm:py-16 px-4 sm:px-6 w-full max-w-4xl">
          <h3 className="text-2xl sm:text-3xl font-bold text-center">Features</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 sm:gap-8 mt-8">
            <div className="p-6 glass text-center">
              <h4 className="text-xl font-semibold">Passwordless Authentication</h4>
              <p className="text-gray-300 mt-2">Leverages passkeys for a seamless and secure login experience.</p>
            </div>
            <div className="p-6 glass text-center">
              <h4 className="text-xl font-semibold">User-Friendly</h4>
              <p className="text-gray-300 mt-2">Easy-to-use interface for seamless experience.</p>
            </div>
            <div className="p-6 glass text-center">
              <h4 className="text-xl font-semibold">Open Source</h4>
              <p className="text-gray-300 mt-2">Completely open-source for transparency and trust.</p>
            </div>
          </div>
        </section>
      </main>

      <footer className="text-center py-4 sm:py-6 glass m-4 px-4">
        <p>&copy; 2025 Pocket ID. All rights reserved.</p>
      </footer>
    </div>
  );
};

export default Home;
