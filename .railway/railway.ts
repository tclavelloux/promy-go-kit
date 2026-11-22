import { defineRailway, project, service } from "railway/iac";

// Starting point for a new Promy Go service: rename the partial and the service to the new Railway service name.
// One Railway project ("main") holds the whole fleet; each repo owns only its own service through a named partial.
// CLI 5.75.0 plans an omitted field as "unset", not "leave alone": once the service exists, every live field must be declared here.
export const partial = "promy-template-go";

export default defineRailway(() =>
  project("main", {
    resources: [
      // No startCommand: the image ENTRYPOINT ["/app"] + CMD ["run"] is the single source.
      service("promy-template-go", {
        // Once the service exists, declare its GitHub source, or apply disconnects it:
        //   source: github("tclavelloux/<repo>", { checkSuites: true }),   (import github; branch defaults to "main")
        build: { builder: "DOCKERFILE", dockerfilePath: "Dockerfile" },
        deploy: {
          // If App Sleeping is on for the service, declare sleepApplication: true, or apply turns it off.
          // Migrations run here, in a separate container, before the deploy; the app does not migrate at boot.
          // A failing command stops the deployment and the previous one keeps serving.
          preDeployCommand: ["/app migrate-database up"],
          preDeployTimeoutSeconds: 300,
          healthcheckPath: "/healthz",
          healthcheckTimeout: 150,
          restartPolicyType: "ON_FAILURE",
          restartPolicyMaxRetries: 10,
        },
        // Variables are managed by this file: declare every Railway variable by name only, values stay in Railway:
        //   env: { NAME: preserve() },   (import preserve)
        // A variable added in the Railway UI must be added here too, or `railway config plan` fails.
      }),
    ],
  }),
);
