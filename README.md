# Senior Software Engineer Test - Answer Template

**Candidate Name:** Zilmas Arjuna Brata Sutrisno
**Position:** Senior Software Engineer 

---

## Part 1: System Design

### Task 1: Architecture Design

#### 1.1 Architecture Design

![Architecture Design Diagram](docs/diagram.png)


**Components:**
- **Ingest Gateway (Logstash):**
  - **Purpose:** Receives, normalizes, and routes security events from Wazuh indexer to downstream processing systems
  - **Inputs:** Wazuh Indexer (wazuh-alerts-* indices).
  - **Outputs:** Event are published to Kafka topic: `event.raw`.
  - **Responsibilities:**
    - Receive the event from wazuh indexer
    - Normalize and parse JSON logs.
    - Forward normalized data to Kafka.
  - **Notes:**
    - Logstash pipeline with JSON and Kafka output plugin.

- **IAM Service:**
  - **Purpose:** Handles authentication and authorization for both users (analysts/admins) and internal services.
  - **Inputs:** Requests from UI, API Gateway
  - **Outputs:** JWT tokens / internal service credentials
  - **Responsibilities:**
    - Manage user accounts, roles
    - Provide /whoami and /verify endpoints for internal service auth.
  - **Notes:** 
    - Go Service
    
- **Rule Service:**
  - **Purpose:** Manage versioned detection rules, suppressions, and policy configurations integrated with Wazuh Manager.
  - **Inputs:**
    - Analyst / Admin via UI or API.
    - Automation proposals (from Scoring / Feedback).
  - **Outputs:**
    - YAML/JSON rule configs.
    - Approved rule changes pushed to Wazuh via REST.
  - **Responsibilities:**
    - CRUD for detection rules, decoders, and suppression policies.
    - Syntax validation and rule testing against sample logs before deployment.
    - Rule versioning, auditing, and rollback tracking.
  - **Notes:**
    - Go Service
    - Audit log every rule push.

- **Alert Intelligence Service:**
  - **Purpose:** Enrich raw Wazuh events with contextual data, then perform short-term correlation and aggregation into alert candidates.
  - **Inputs:** Kafka topic: `event.raw`
  - **Outputs:** Kafka topic: `alert.intelligence` (enriched + correlated events)
  - **Responsibilities:**
    - Add asset metadata (owner, critically, tags, maintenance window)
    - Add threat intel (IP/domain reputation, ASN, geo).
    - Add vulnerability info (CVE information).
    - Correlation aggregate repeated alerts from the same IP to the same target
    - Correlation behaviour between event

  - **Notes:** 
    - Go Service
    - Use redis for enrichment (asset metadata, threat intel, cve)

- **Scoring Service:**
  - **Purpose:** Use ML/AI models to assign risk scores and decisions to enriched/correlated alerts.
  - **Inputs:** Kafka topic: `alert.intelligence`
  - **Outputs:** Kafka topic: `alert.scored`
  - **Responsibilities:**
    - Extract alert features and run inference on a trained model.
    - Produce {score, confidence, reasons, recommendation}.
    - Tag “auto-suppress”, “needs-review”, or “critical-escalate”.
  - **Notes:**
    - Python (fast api)

- **Feedback Service:**
  - **Purpose:** Collect analyst validation labels (True Positive, False Positive, Needs Info) and teach the AI model how to improve its scoring accuracy over time.
  - **Inputs:** 
    - Analyst actions from UI.
    - Scored alerts from Kafka `alert.scored`.
  - **Outputs:**
    - Feedback events `alert.feedback`
  - **Responsibilities:**
    - Store immutable feedback record (alert_id, analyst, label, reason).
    - Detect recurring False Positive patterns and propose suppression rules.
    - Trigger retraining pipeline updates for the AI model.
  - **Notes:** 
    - Python (fast api)

- **SOAR Service:**
  - **Purpose:** Execute automated playbooks (containment, notification, ticketing) based on scored alert decisions.
  - **Inputs:**
    - Kafka topic: `alert.scored`.
    - Manual trigger from UI.
  - **Outputs:**
    - action logs
    - Integration with external systems (email, Slack, firewall API).
  - **Responsibilities:**
    - Parse recommendation: escalate/contain/suppress.
    - Trigger response flows
    - Record all playbook executions and results.
    - Action logging and audit trails
    - Support rollback actions and simulation/dry-run mode.
  - **Notes:**
    - Go Service
    - Log Every Action

- **API Gateway:**
  - **Purpose:** Unified secure entry point for UI, third-party integrations, and inter-service calls.
  - **Inputs:** HTTP/gRPC requests from UI and clients.
  - **Outputs:** Proxy routes to internal services (IAM, Rule, Feedback, SOAR, Scoring results).
  - **Responsibilities:**
    - JWT validation via IAM Service.
    - Handle routing, rate limiting, and access logging.
  - **Notes:**
    - Can use Traefik, API Six, or NGINX.

- **UI:**
  - **Purpose:** Provide SOC analysts and admins a unified dashboard for alert review, feedback submission, rule management, and SOAR actions.
  - **Inputs:**
    - REST API via API Gateway.
    - WebSocket/Server-Sent Events for real-time alerts.
  - **Outputs:** Visualized analyst feedback, rule proposals, SOAR triggers.
  - **Responsibilities:**
    - Visualized analyst feedback, rule proposals, SOAR triggers.
  - **Notes:**
    - Vite

**Feedback Loop from analyst validation improve future precision**
1. Analyst Validation (via UI):
  - Analyst mark alert as `True Positive` or `False Positive` and there reason as well on the dashboard
  - the labeled data then sent to the `Feedback Service`.
2. Feedback Processing:
  - `Feedback Service` store the labeled event and forward this data to the retraining pipeline. If recurring False Positive patterns are detected, the service sends a suppression proposal to the `Rule Service`.
3. Model & Rule Improvement
  - `Scoring Service` uses the newly retrained model to improve the accuracy of future alert scoring.
  - `Rule Service` applies approved suppression or adjustment rules, directly reducing false positives in Wazuh.
4. Result
  Over time, the system continuously becomes more precise — reducing alert noise, improving response speed, and lowering the false positive rate.

#### 1.2 Data Flow Explanation
**1. Agent Log**
- Each `wazuh-agent` continuously gathers logs from hosts, servers, network devices.
- `wazuh-agent` buffer logs locally before sending them to `wazuh-manager`.
- `wazuh-agent` forward the raw logs to `wazuh-manager`.

**2. Wazuh Manager**
- `wazuh-manager` receive the logs and decode and normalize the event using `decoders`.
- after that applies rules matching logic to detect anomalies or suspious activity using `rules`
- after match with that rules will forward to `wazuh-indexer` under `wazuh-alert-*`

**3. Elastic Indexer -> Ingest Gateway**
- The `Ingest Gateway (Logstash)` subscribes to the Elastic indices. , extracts new alerts, and normalizes the payload into a common JSON schema.
- These normalized alerts are pushed into Kafka topic `event.raw`.

**4. Alert Intelligence Service**
- Consumes `event.raw` and performs enrichment + correlation:
  - Asset enrichment: add asset metadata (owner, critically, tags, maintenance window)
  - Threat intelligence: threat intel (IP/domain reputation, ASN, geo) and vulnerability info (CVE information).
  - Correlation: aggregate repeated alerts from the same IP to the same target and behaviour between event.
- Outputs are publishes to Kafka topic `alert.intelligence` for downstream AI scoring.

**5. Scoring Service**
- Consume `alert.intelligence`.
- Extracts contextual features and sends them inference ML/AI.
- Produces a risk score (0–1) and a recommendation (auto-suppress, needs-review, critical-escalate).
- Output publishes to Kafka topic alert.scored.
 

**6. Analyst Review & Feedback**
- The UI receives scored alerts (through the API Gateway).
- Analysts review only relevant alerts:
  - High-score → escalated to SOAR.
  - Medium-score → manual review.
  - Low-score → auto-suppressed (noise bucket).
- Analysts submit feedback (TP/FP) through the UI → `Feedback Service`, closing the loop for model retraining and rule suppression proposals.

#### 1.3 Scalability & Fault Tolerance
**Distribute event processing across nodes**
**1. Wazuh Agent** 
Enable local buffering: queue_size tuned per host (avoid drops on wazuh server).
**2. Logstash (Ingest Gateway) → Kafka**
- Logstash persistent queue enabled.
- Kafka Settings:
  - `acks=all` for guaranted no data loss
  - `enable.idempotence=true`  for guaranted exactly-once delivery semantics
**3. Kafka -> Microservices**
- Consumer Group. each service scales horizontally by add adding consumer in the same group.

**Handle API rate limits & node failures gracefully**
**1. Exponential backoff with jitter**
- Add jitter (a small random offset) to avoid all clients retrying at the same time. to prevents cascading overload when APIs are unstable.
```
200 ms → 400 ms → 800 ms → 1.6 s → 3.2 s … up to 30 s
```

**2. Caching**
- Using redis if there is api that the data does'nt frequently change like config, asset or cve.

**Maintain system state consistency**
**1. Shadow state**
- Rule Service stores the intended rules/decoders/suppressions in Git
- Every change is a Change Set (ID, author, rollout plan)
**2. Hourly reconcile job**
Pull active rules from Wazuh Manager (REST/SSH read).
**3. Idempotency & exactly-once effects**
All policy pushes are idempotent (same request → same result). Include content hash in the API payload.

#### 1.4 Ethical & Operational Constraints 
**1. Confidence Thresholding**
- The Scoring Service only auto-suppresses alerts when the AI confidence is very low (e.g., < 0.3).
- Alerts with medium or uncertain confidence go to manual analyst review.
- High-risk alerts are always escalated to SOAR or human analysts.

**2. Immutable Event Trail**
- Every action—rule change, model update, SOAR command—is write on audit db.

**3. Rollback & Version Control**
- Rule Service and SOAR maintain Git-style versioning:
  - Every rule or playbook change = new commit with diff & approver.
  - Rollback = one-click restore of previous version.

**4. Access Control & RBAC**
- There is role for review and feeback and approval for recheck. Separates duties and prevents misuse of automation. 


## Part 2: Coding Challenge
**Complete Implementation Reference**: Please see the `automation-wazuh-triage/` folder for the full working solution.

This folder contains a comprehensive Wazuh Security Event Triage Automation System with Clean Architecture implementation, including all source code, API documentation, and usage examples.

## Part 3: Analytical Case Study

### Case Study Overview
```
[Summarize the case study scenario]
```

### Problem Analysis

#### 3.1 Current State Assessment
- [ ] **Strengths:** [What's working well]
- [ ] **Weaknesses:** [Pain points and bottlenecks]
- [ ] **Opportunities:** [Areas for improvement]
- [ ] **Threats:** [Potential risks]

#### 3.2 Root Cause Analysis
1. **Primary Issues:**
   - [ ] Issue 1: [Description and impact]
   - [ ] Issue 2: [Description and impact]

2. **Contributing Factors:**
   - [ ] [Technical factors]
   - [ ] [Process factors]
   - [ ] [Resource factors]

### Proposed Solutions

#### 3.3 Short-term Solutions (0-3 months)
- [ ] **Solution 1:** [Description, effort, impact]
- [ ] **Solution 2:** [Description, effort, impact]

#### 3.4 Medium-term Solutions (3-12 months)
- [ ] **Solution 1:** [Description, effort, impact]
- [ ] **Solution 2:** [Description, effort, impact]

#### 3.5 Long-term Solutions (1+ years)
- [ ] **Solution 1:** [Description, effort, impact]
- [ ] **Solution 2:** [Description, effort, impact]

### Implementation Plan
```
[Timeline, resources, milestones, success metrics]
```

### Risk Assessment
- [ ] **Risk 1:** [Description, probability, impact, mitigation]
- [ ] **Risk 2:** [Description, probability, impact, mitigation]

---

## Part 4: Behavioral & Design Reasoning

### 4.1 Leadership Experience

#### Question: [Insert specific behavioral question]
**Situation:** [Describe the context]  
**Task:** [What needed to be accomplished]  
**Action:** [What you did specifically]  
**Result:** [Outcome and impact]  

### 4.2 Technical Decision Making

#### Question: [Insert specific question about technical decisions]
**Challenge:** [Describe the technical challenge]  
**Options Considered:** 
- [ ] Option 1: [Pros/Cons]
- [ ] Option 2: [Pros/Cons]
- [ ] Option 3: [Pros/Cons]

**Decision Made:** [Your choice and reasoning]  
**Outcome:** [Results and lessons learned]  

### 4.3 Design Philosophy

#### Question: [Insert design-related question]
**Design Principles Applied:**
- [ ] [Principle 1 and application]
- [ ] [Principle 2 and application]
- [ ] [Principle 3 and application]

**Trade-offs Considered:**
- [ ] [Trade-off 1: Benefits vs. Costs]
- [ ] [Trade-off 2: Benefits vs. Costs]

### 4.4 Problem-Solving Approach

#### Question: [Insert problem-solving question]
**Problem Definition:** [How you understood the problem]  
**Investigation Process:** [Your approach to gathering information]  
**Solution Development:** [How you developed solutions]  
**Implementation:** [Execution strategy]  
**Validation:** [How you verified success]  

### 4.5 Collaboration & Communication

#### Question: [Insert collaboration question]
**Context:** [Team/project situation]  
**Challenge:** [Communication or collaboration obstacle]  
**Approach:** [Your strategy for working with others]  
**Outcome:** [Results and relationship impact]  

---

## Additional Considerations

### Assumptions Made
- [ ] [List any assumptions you made while solving the problems]

### Questions for Clarification
- [ ] [Questions you would ask in a real-world scenario]

### Future Enhancements
- [ ] [Ideas for extending or improving the solutions]

---

**Note:** This template is designed to be comprehensive. Fill in each section with detailed, specific responses that demonstrate your technical expertise, problem-solving abilities, and leadership experience.
