import fs from "node:fs/promises";
import path from "node:path";

const root = process.cwd();
const aofRoot = path.join(root, ".aof");
const visibilityDir = path.join(aofRoot, "artifacts", "visibility");

async function readJsonIfPresent(filePath) {
  try {
    const raw = await fs.readFile(filePath, "utf8");
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

async function readLatestJsonFromDir(dirPath) {
  try {
    const entries = await fs.readdir(dirPath);
    const jsonFiles = entries.filter((entry) => entry.endsWith(".json")).sort();
    if (jsonFiles.length === 0) {
      return null;
    }
    return readJsonIfPresent(path.join(dirPath, jsonFiles.at(-1)));
  } catch {
    return null;
  }
}

async function readLatestTextFromDir(dirPath) {
  try {
    const entries = await fs.readdir(dirPath);
    const files = entries.filter((entry) => entry.endsWith(".md")).sort();
    if (files.length === 0) {
      return null;
    }
    return await fs.readFile(path.join(dirPath, files.at(-1)), "utf8");
  } catch {
    return null;
  }
}

function normalizeGoalContent(goalProjection, fallback) {
  return goalProjection?.content ?? fallback ?? null;
}

function buildStatusCard({
  session,
  operatingGoal,
  nextValueSlice,
  latestDecisionPath,
  cadenceSchedule,
  cadenceProfile
}) {
  return {
    view_type: "status_card",
    as_of: new Date().toISOString(),
    usage_level: "runtime mandatory",
    current_phase: session?.current_stage ?? session?.status ?? "bootstrap",
    current_goal: normalizeGoalContent(operatingGoal, "Bootstrap the GCP learning repository"),
    owner: "Facilitator",
    open_signals: [],
    next_checkpoint: "Review latest council output and advance the next value slice",
    latest_artifact_ref: latestDecisionPath ?? "none",
    runtime_evidence_state: session ? "present" : "partial",
    north_star_goal: "Build a reusable GCP practical learning repository for PM/SE onboarding",
    current_operating_goal: normalizeGoalContent(
      operatingGoal,
      "Turn the repository into a runtime-backed AOF learning project"
    ),
    next_value_slice: normalizeGoalContent(
      nextValueSlice,
      "Run the AOF clarification and planning flow for the current repository state"
    ),
    open_task_count: 1,
    stale_task_count: 0,
    recent_confirmation_ref: "context/active/recent-confirmation-window.json",
    cadence_timing_state: cadenceSchedule?.timing_state ?? "not-started",
    cadence_scheduler_state: cadenceSchedule?.scheduler_state ?? "manual",
    cadence_dispatch_state: cadenceSchedule?.recommended_command ? "ready" : "manual",
    cadence_scheduler_profile: cadenceProfile?.selected_profile ?? "manual",
    cadence_next_check_at: cadenceSchedule?.recommended_next_check_at ?? null,
    cadence_reason: cadenceSchedule?.reason ?? "Visibility is currently refreshed manually from live .aof state"
  };
}

function buildTimelineFeed({ session, latestDecisionText }) {
  const entries = [];
  if (session) {
    entries.push({
      at: session.updated_at ?? session.created_at ?? new Date().toISOString(),
      actor: "Facilitator",
      event_type: session.status ?? "session_created",
      summary: `AOF session ${session.session_id} is active for this repository`,
      rationale: "Runtime-backed project tracking is now enabled for the learning repository",
      next: "Advance council planning and update the next value slice",
      refs: [session.session_id]
    });
  }
  if (latestDecisionText) {
    entries.push({
      at: new Date().toISOString(),
      actor: "Council of Three",
      event_type: "decision_record",
      summary: "Decision record exists for the bootstrap learning slice",
      rationale: "The repository keeps explicit artifact and outcome intent in AOF form",
      next: "Continue with runtime-managed clarification and approval history",
      refs: ["DEC-001-bootstrap-learning-repo"]
    });
  }
  return {
    view_type: "timeline_feed",
    entries
  };
}

function buildFlowSnapshot({ session }) {
  const currentStage = session?.current_stage ?? session?.status ?? "clarification";
  const current =
    currentStage === "planning" || currentStage === "approval" || currentStage === "proposal"
      ? "planning"
      : "clarification";
  const nodes = [
    { id: "request", label: "request_received", state: "done" },
    { id: "clarification", label: "clarification", state: current === "clarification" ? "current" : "done" },
    {
      id: "planning",
      label: "planning",
      state: currentStage === "approval" ? "done" : current === "planning" ? "current" : "pending"
    },
    { id: "approval", label: "approval", state: currentStage === "approval" ? "current" : "pending" },
    { id: "delivery", label: "delivery", state: "pending" }
  ];
  return {
    view_type: "flow_snapshot",
    nodes,
    edges: [
      { from: "request", to: "clarification", reason: "initial framing started" },
      { from: "clarification", to: "planning", reason: "clarification answers accepted" },
      { from: "planning", to: "approval", reason: "council planning completed" },
      { from: "approval", to: "delivery", reason: "approved slice is implemented" }
    ],
    current_node: current,
    open_branches: []
  };
}

async function main() {
  await fs.mkdir(visibilityDir, { recursive: true });

  const session = await readLatestJsonFromDir(path.join(aofRoot, "sessions"));
  const operatingGoal = await readJsonIfPresent(path.join(aofRoot, "goals", "operating-goal.json"));
  const nextValueSlice = await readJsonIfPresent(path.join(aofRoot, "goals", "next-value-slice.json"));
  const cadenceSchedule = await readJsonIfPresent(path.join(aofRoot, "context", "active", "cadence-schedule.json"));
  const cadenceProfile = await readJsonIfPresent(path.join(aofRoot, "context", "active", "cadence-scheduler-profile.json"));
  const latestDecisionText = await readLatestTextFromDir(path.join(aofRoot, "decisions"));

  const latestDecisionPath = latestDecisionText ? ".aof/decisions/latest" : null;
  const statusCard = buildStatusCard({
    session,
    operatingGoal,
    nextValueSlice,
    latestDecisionPath,
    cadenceSchedule,
    cadenceProfile
  });
  const timelineFeed = buildTimelineFeed({ session, latestDecisionText });
  const flowSnapshot = buildFlowSnapshot({ session });

  await Promise.all([
    fs.writeFile(path.join(visibilityDir, "status-card.json"), JSON.stringify(statusCard, null, 2)),
    fs.writeFile(path.join(visibilityDir, "timeline-feed.json"), JSON.stringify(timelineFeed, null, 2)),
    fs.writeFile(path.join(visibilityDir, "flow-snapshot.json"), JSON.stringify(flowSnapshot, null, 2))
  ]);

  process.stdout.write(
    JSON.stringify(
      {
        ok: true,
        status: path.join(".aof", "artifacts", "visibility", "status-card.json"),
        timeline: path.join(".aof", "artifacts", "visibility", "timeline-feed.json"),
        flow: path.join(".aof", "artifacts", "visibility", "flow-snapshot.json")
      },
      null,
      2
    )
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
