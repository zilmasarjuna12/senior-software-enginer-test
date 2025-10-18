# Senior Software Engineer Test - Answer Template

**Candidate Name:** Zilmas Arjuna Brata Sutrisno
**Position:** Senior Software Engineer 

---

## Part 1: System Design

### Task 1: Architecture Design

#### 1.1 High-Level Architecture Overview

![System Architecture Diagram](docs/diagram.png)


**Components:**
- **Ingest Gateway (Logstash):**
  - **Purpose:** Receives, normalizes, and routes security events from Wazuh indexer to downstream processing systems
  - **Inputs:** 
    - Wazuh Indexer (wazuh-alerts-* indices)
    - Optional other sensors (Suricata, Zeek, OT telemetry).
  - **Outputs:** Kafka topic: event.raw
  - **Responsibilities:**
    - Receive the event from wazuh indexer
    - Forward to Kafka (event.raw).
  - **Notes:**
    - Logstash pipeline with JSON codec and Kafka output plugin.

- **IAM Service:**
  - **Purpose:** Provide authentication, authorization, and tenant/user context for analysts, API clients, and internal services.
  - **Inputs:** Requests from UI, API Gateway
  - **Outputs:** JWT tokens / service credentials
  - **Responsibilities:**
    - Manage user accounts, roles
    - Provide /whoami and /verify endpoints for internal service auth.
  - **Notes:** 
    - Centralized identity for SOC microservices.
    
- **Rule Service:**
  - **Purpose:** Manage versioned detection rules, suppressions, and policy configurations integrated with Wazuh Manager.
  - **Inputs:**
    - Analyst / Admin via UI or API.
    - Automation proposals (from Scoring / Feedback).
  - **Outputs:**
    - YAML/JSON rule configs.
    - Approved rule changes pushed to Wazuh via REST.
  - **Responsibilities:**
    - CRUD for custom rules/decoders/suppressions.
    - Validate syntax and test new rules against sample logs before deployment.
    - Rule versioning and deployment
    - Pattern matching and correlation
    - Performance optimization and caching
  - **Notes:**
    - Go Service
    - Audit log every rule push.

- **Alert Intelligence Service:**
  - **Purpose:** Enrich raw Wazuh events with contextual data, then perform short-term correlation and aggregation into alert candidates.
  - **Inputs:** Kafka topic: event.raw
  - **Outputs:** Kafka topic: alert.intelligence (enriched + correlated events)
  - **Responsibilities:**
    - Add asset metadata (owner, critically, tags, maintenance window)
    - Add threat intel (IP/domain reputation, ASN, geo).
    - Add vulnerability info (CVE status).
    - Correlation aggregate repeated alerts from the same IP to the same target
    - Correlation behaviour between event

  - **Notes:** 
    - Go Service
    - Use redis for enrichment (asset metadata, threat intel, cve)

- **Scoring Service:**
  - **Purpose:** Use ML/AI models to assign risk scores and decisions to enriched/correlated alerts.
  - **Inputs:** Kafka topic: alert.intelligence
  - **Outputs:** Kafka topic: alert.scored
  - **Responsibilities:**
    - Extract features from alert payloads intelegence
    - Call ML inference
    - Produce {score, confidence, reasons, recommendation}.
    - Tag “auto-suppress”, “needs-review”, or “critical-escalate”.
  - **Notes:**
    - Python (fast api)

- **Feedback Service:**
  - **Purpose:**Collect analyst validation labels (True Positive, False Positive, Needs Info) and teach the AI model how to improve its scoring accuracy over time.
  - **Inputs:** 
    - Analyst actions from UI.
    - Scored alerts from Kafka (alert.scored).
  - **Outputs:**
    - Feedback events (alert.feedback Kafka).
  - **Responsibilities:**
    - Store immutable feedback record (alert_id, analyst, label, reason).
    - Trigger retraingin pipeline.
    - Allow rule suppression proposal generation for consistent FP patterns.
  - **Notes:** 
    - Python (fast api)

- **SOAR Service:**
  - **Purpose:** Execute automated playbooks (containment, notification, ticketing) based on scored alert decisions.
  - **Inputs:**
    - Kafka topic: alert.scored.
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
    - JWT validation via IAM.
    - Rate limiting and logging.
    - routing.
  - **Notes:**
    - Traefik / ApiSix NGInx.
    - Enforce RBAC between services.

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

**Technology Stack:**
- **Frontend:** [Technology choices and rationale]
- **Backend:** [Technology choices and rationale]
- **Database:** [Technology choices and rationale]
- **Infrastructure:** [Technology choices and rationale]
- **Message Queue/Event Bus:** [If applicable]
- **Caching Layer:** [If applicable]

#### 1.2 Architecture Patterns
- **Pattern Used:** [e.g., Microservices, Monolith, Event-Driven, Layered, etc.]
- **Justification:** [Why this pattern was chosen for this specific use case]
- **Trade-offs:** [What you gained vs. what you sacrificed]

#### 1.3 API Design
```yaml
# Example API specification
/api/v1/users:
  GET: [Description]
  POST: [Description]
  
/api/v1/users/{id}:
  GET: [Description]
  PUT: [Description]
  DELETE: [Description]
```

#### 1.4 Database Schema
```sql
-- [Database schema design]
-- Include tables, relationships, and key constraints
```

#### 1.5 Security Considerations
- [ ] **Authentication:** [Method and implementation]
- [ ] **Authorization:** [RBAC, permissions model]
- [ ] **Data Protection:** [Encryption, PII handling]
- [ ] **API Security:** [Rate limiting, input validation]

### Task 2: Data Flow Explanation

#### 2.1 Request Flow Diagram
```
[User] → [Load Balancer] → [API Gateway] → [Service Layer] → [Database]
```

#### 2.2 Detailed Data Flow Steps

**Step 1: User Request**
- [ ] [Description of initial user action]
- [ ] [Input validation process]

**Step 2: Processing Layer**
- [ ] [How the system processes the request]
- [ ] [Business logic implementation]

**Step 3: Data Persistence**
- [ ] [How data is stored/retrieved]
- [ ] [Transaction handling]

**Step 4: Response Generation**
- [ ] [How response is formatted]
- [ ] [Error handling mechanism]

#### 2.3 Data Synchronization
- [ ] **Strategy:** [Eventual consistency, Strong consistency, etc.]
- [ ] **Implementation:** [How data sync is handled across services]

---

### Task 3: Scalability & Fault Tolerance

#### 3.1 Scalability Strategy

**Horizontal Scaling:**
- [ ] [Auto-scaling policies]
- [ ] [Load balancing strategy]
- [ ] [Database sharding/partitioning]

**Vertical Scaling:**
- [ ] [Resource optimization]
- [ ] [Performance tuning]

**Caching Strategy:**
- [ ] **Application Level:** [Redis, Memcached, etc.]
- [ ] **Database Level:** [Query optimization, indexing]
- [ ] **CDN:** [Static content delivery]

#### 3.2 Fault Tolerance Mechanisms

**Redundancy:**
- [ ] [Multi-region deployment]
- [ ] [Database replication]
- [ ] [Service redundancy]

**Circuit Breaker Pattern:**
- [ ] [Implementation details]
- [ ] [Fallback mechanisms]

**Monitoring & Alerting:**
- [ ] [Health checks]
- [ ] [Performance monitoring]
- [ ] [Error tracking]

**Disaster Recovery:**
- [ ] [Backup strategy]
- [ ] [Recovery procedures]
- [ ] [RTO/RPO targets]

---

## Part 2: Coding Challenge

### Problem Statement
```
[Copy the exact problem statement here]
```

### Solution Approach
```
[Explain your approach and algorithm choice]
```

### Implementation

#### Language: [Programming Language Used]

```[language]
// [Your code implementation here]
// Include comments explaining key logic
```

### Time & Space Complexity
- **Time Complexity:** O([complexity])
- **Space Complexity:** O([complexity])
- **Justification:** [Explain the complexity analysis]

### Test Cases
```[language]
// [Test cases covering edge cases]
```

### Alternative Solutions
```
[Discuss alternative approaches and trade-offs]
```

---

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
