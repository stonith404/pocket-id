import "@fortawesome/fontawesome-free/css/all.min.css";
import React from "react";
import FeatureBox from "../components/feature-box";
import "/styles.css";

const Home: React.FC = () => {
  return (
    <div className="text-white h-screen flex flex-col bg-muted/40">
      <header style={{ backgroundColor: "hsl(240, 10%, 3.9%)" }}>
        <div className="w-full border-b border-black">
          <div className="container flex w-full items-center justify-between px-4 md:px-10">
            <div className="flex h-16 items-center">
              <img
                src="https://docs.pocket-id.org/img/pocket-id.png"
                alt="Pocket ID Logo"
                className="mr-3 h-8 w-8"
              />
              <h2 className="text-sm font-medium" style={{ margin: 0 }}>
                Pocket ID
              </h2>
            </div>
            <a
              href="https://github.com/stonith404/pocket-id"
              target="_blank"
              rel="noopener noreferrer"
              style={{ color: "hsl(0, 0%, 98%)" }}
              className="text-white text-2xl"
            >
              <i className="fab fa-github" aria-hidden="true"></i>
            </a>
          </div>
        </div>
      </header>

      <main className="flex-1 flex flex-col justify-center items-center px-4 sm:px-0 container">
        <section className="flex items-center mt-10 flex-col-reverse lg:flex-row gap-5">
          <div>
            <h1 className="text-sm font-extrabold">
              Secure Your Services with OIDC
            </h1>
            <p className="mt-4 text-lg">
              Pocket ID is a simple and easy-to-use OIDC provider that allows
              users to authenticate with their passkeys to your services.
            </p>
            <a
              href="/docs/introduction"
              className="mt-6 inline-block text-black px-6 py-3 rounded-lg font-semibold"
              style={{
                backgroundColor: "hsl(0, 0%, 98%)",
                color: "hsl(240, 10%, 3.9%)",
              }}
            >
              Get Started
            </a>
          </div>
          <img
            src="/img/landing/authorize-screenshot.png"
            alt="Pocket ID Logo"
            className="max-h-[350px] xl:max-h-[450px]"
          />
        </section>

        <section className="mt-15">
          <h2 className="!text-3xl font-bold">Features</h2>
          <div className="flex flex-col gap-5">
            <FeatureBox
              title="Passwordless Authentication"
              description="Pocket ID only supports passkey authentication, which means you don't need a password."
              imgSrc="/img/landing/passkey-auth-screenshot.png"
            />
            <FeatureBox
              title="Restrict User Groups"
              description="You can select which user groups are allowed to authenticate with your services."
              imgSrc="/img/landing/allowed-usergroups-screenshot.png"
              imgLeft={false}
            />
            <FeatureBox
              title="Audit Logs"
              description="Keep track of your account activities. If SMTP is configured, you can even receive sign-in notifications."
              imgSrc="/img/landing/audit-log-screenshot.png"
            />
            <FeatureBox
              title="LDAP"
              description="Sync your users and groups from your LDAP server to Pocket ID."
              imgSrc="/img/landing/ldap-screenshot.png"
              imgLeft={false}
            />
          </div>

          <p className="!mt-5 text-center">And much more...</p>
        </section>
      </main>

      <div className="flex flex-col items-center mt-10">
        <p className="py-3 text-xs text-muted-foreground">
          &copy; 2025 Pocket ID
        </p>
      </div>
    </div>
  );
};

export default Home;
