## ✅ Summary: Running Steampipe on ECS with `credential_source = EcsContainer`

## ✅ Prerequisites: Steampipe on ECS with `credential_source = EcsContainer`

Before you start, ensure the following:

---

### 🔹 **1. ECS Fargate Setup**

- An ECS **cluster** is already created.
- A **task definition** exists with:
  - `FARGATE` launch type.
  - A **task role** with permissions to assume other roles (`sts:AssumeRole`).
  - A lightweight image (e.g., `amazonlinux:2`).
- A **service** or **standalone task** is running from this task definition.
- ECS Exec is **enabled** on the service or task (`EnableExecuteCommand: true`).

---

### 🔹 **2. IAM Role Requirements**

#### **Task Role (Assigned to the ECS Task):**
- Must include:
  - `sts:AssumeRole`
  - `ssmmessages:*` permissions (for ECS Exec)
- Trust relationship must allow `ecs-tasks.amazonaws.com`.

#### **Target Role (To be assumed using credential_source):**
- Example: `OrganizationAccountAccessRole` or `test_admin`
- Must allow the task role to assume it via trust policy.

---

### 🔹 **3. Networking Requirements**

- The task must be launched in **public or private subnets** with:
  - **NAT Gateway or Internet Access** (to pull Steampipe & plugins).
  - Security group that allows **outbound HTTPS (443)** traffic.
- VPC endpoints for Systems Manager (`ssm`, `ssmmessages`, `ec2messages`) if in private subnets.

---

### 🔹 **4. System Manager Setup (for ECS Exec)**

- Systems Manager agent is **enabled by default** in Fargate.
- No extra agent setup is needed, but:
  - Your account must have **SSM Session Manager plugin** installed locally.
  - The task must run in a supported platform version (Fargate **1.4.0 or newer**).

---

### 🔹 **5. Local Requirements (for `aws ecs execute-command`)**

- You must have the following installed **locally** (on your machine):
  - [AWS CLI v2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
  - [Session Manager Plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)
  - AWS credentials with permission to use ECS Exec

---

### 🔹 **6. Steampipe Requirements (Installed Inside the Container)**

Inside the ECS container, you will install:

- `tar`
- `gzip`
- `unzip`
- `which`
- `shadow-utils` (for `useradd`)
- `util-linux` (for `su`)

These are required for:
- Creating a non-root user
- Installing Steampipe
- Switching to the new user for plugin execution

---

NOTE: You could see the sample [CloudFormation stack](https:github.com/ParthaI/linear-utils/add-initial-commits/docs/run-steampipe-in-ecs-container/sample-cloudfromation-stack/cf-stack.yaml) to do all the setup for you.

### **1. Start an interactive shell into your ECS container**
Using `aws ecs execute-command`:

```bash
aws ecs execute-command \
  --cluster <cluster-name> \
  --task <task-id> \
  --container <container name> \
  --command "/bin/sh" \
  --interactive
```

This logs you into the container as the **root user**.

---

### **2. Install required dependencies**
Inside the container, install necessary packages:

```bash
yum install -y tar gzip unzip which shadow-utils util-linux
```

---

### **3. Install Steampipe**
As root, run:

```bash
curl -fsSL https://steampipe.io/install/steampipe.sh | sh
```

> ⚠️ Steampipe cannot run as root, so this just sets up the binary.

---

### **4. Create and switch to a non-root user**

Still as root:

```bash
useradd -m steampipe
echo $AWS_CONTAINER_CREDENTIALS_RELATIVE_URI  # save the value
```

Then:

```bash
su - steampipe
```

Inside the `steampipe` user shell:

```bash
export AWS_CONTAINER_CREDENTIALS_RELATIVE_URI=/v2/credentials/<value-from-root>
```

This gives the new user access to the ECS credentials.

---

### **5. Configure AWS CLI profile using ECS credentials**

```bash
account_id=$(aws sts get-caller-identity --query Account --output text)

mkdir -p ~/.aws
cat <<EOF >> ~/.aws/config
[profile aws_${account_id}]
role_arn = arn:aws:iam::${account_id}:role/test_admin
credential_source = EcsContainer
EOF
```

Test it:

```bash
aws sts get-caller-identity --profile aws_${account_id}
```

---

### **6. Install and configure Steampipe AWS plugin**

```bash
steampipe plugin install aws
```

Create your connection config at `~/.steampipe/config/aws.spc`:

```hcl
connection "aws" {
  plugin  = "aws"
  regions = ["*"]
  profile = "aws_${account_id}"
}
```

---

### **7. Run your Steampipe queries**

```bash
steampipe query
> select * from aws_account;
```

✅ You should now be able to query AWS using the assumed role, backed by the ECS container credentials.