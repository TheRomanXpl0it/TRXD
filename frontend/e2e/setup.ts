import { exec } from "child_process";
import { promisify } from "util";
import * as path from "path";

const execAsync = promisify(exec);

export default async function globalSetup() {
  const echoServerPath = path.resolve("../tests/echo-server");

  try {
    const { stdout, stderr } = await execAsync(
      "docker build -t echo-server:latest .",
      {
        cwd: echoServerPath,
        timeout: 60000,
      },
    );

    if (stdout) {
      console.log("Docker build output:", stdout.trim());
    }
    if (stderr) {
      console.log("Docker build stderr:", stderr.trim());
    }

    console.log("Echo-server image built");
  } catch (error: any) {
    console.error("Failed to build echo-server image:", error.message);
    if (error.stdout) console.error("stdout:", error.stdout);
    if (error.stderr) console.error("stderr:", error.stderr);
    throw error;
  }

  console.log("Global setup complete");
}
