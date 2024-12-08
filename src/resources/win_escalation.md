# That's a list of some windows api that are more useful for a privilage escalation
---

### 1. **Impersonation & Token Manipulation (High chance of privilege escalation)**

- `advapi32 ImpersonateLoggedOnUser` – Lets you impersonate a logged-in user, running code with their privileges (not a sandbox).
- `advapi32 ImpersonateAnonymous` – Impersonates an anonymous user, useful for some escalation attacks.
- `advapi32 SetThreadToken` – Set a token for a thread, allowing code to run with the impersonated user's privileges.
- `advapi32 DuplicateToken` – Duplicates an existing token, giving you a new token with privileges.
- `kernel32 CreateProcessWithToken` – Creates a process with a token, useful to run a process with higher privileges.
- `kernel32 CreateProcessAsUser` – Lets you run a process as another user with their privileges.
- `advapi32 RevertToSelf` – Reverts to the original user's privileges, often used after impersonating.
- `advapi32 SetTokenInformation` – Modifies token info, used to tweak process privileges.
- `advapi32 AdjustTokenPrivileges` – Alters privileges of an existing token.
- `advapi32 GetTokenInformation` – Pulls token information, useful to figure out the user’s privileges.
- `kernel32 OpenThreadToken` – Opens a thread's token, useful for accessing elevated privileges.
- `kernel32 OpenThreadTokenEx` – Advanced version of `OpenThreadToken`, gives more granular access to tokens.

### 2. **Security & Access Control (Privileges related to security and permissions)**

- `advapi32 SetSecurityInfo` – Changes security info on an object, useful for modifying file or resource permissions.
- `advapi32 SetNamedSecurityInfo` – Same as `SetSecurityInfo`, but for specific resources.
- `advapi32 AddAccessAllowedAce` – Adds an ACE (Access Control Entry) to allow access to an object.
- `advapi32 AddAccessDeniedAce` – Adds an ACE to deny access to an object.
- `advapi32 SetWindowsStationUser` – Sets a user for a window station, useful for manipulating window sessions.
- `advapi32 SetDesktopUser` – Sets the user for a desktop session, also useful for escalation.
- `advapi32 LookupPrivilegeValue` – Gets the value of a security privilege, useful for privilege escalation attacks.

### 3. **Process & Thread Manipulation (Medium chance of escalation)**

- `kernel32 CreateFileMapping` – Creates a file mapping, can be used for shared memory access.
- `kernel32 OpenThread` – Opens an existing thread, useful for manipulating other threads with elevated privileges.
- `kernel32 SetThreadContext` – Sets a thread’s context, often used in exploits to control process execution.
- `kernel32 SetThreadPriority` – Changes a thread’s priority, can raise privileges when combined with other techniques.
- `kernel32 SetThreadAffinityMask` – Sets a thread’s affinity mask, useful for certain types of attacks.
- `kernel32 SetThreadExecutionState` – Alters a thread’s execution state.
- `kernel32 GetThreadContext` – Gets a thread’s context, useful for understanding its state and manipulating execution.
- `kernel32 GetThreadPriority` – Retrieves a thread’s priority, useful for attacks.

### 4. **Job Object & Mutex (Low chance of direct escalation)**

- `kernel32 CreateProcessWithJobObject` – Creates a process with an associated job object, useful for limiting or monitoring processes.
- `kernel32 AssignProcessToJobObject` – Assigns a process to a job object, can manipulate control over it.
- `kernel32 SetInformationJobObject` – Sets info on a job object, less common for escalation.
- `kernel32 GetInformationJobObject` – Retrieves info on a job object.
- `kernel32 CreateMutex` – Creates a mutex, useful for controlling concurrency between processes.
- `kernel32 CreateEvent` – Creates a synchronization event, less useful for direct escalation.
- `kernel32 CreateSemaphore` – Creates a semaphore for synchronization.
- `kernel32 CreateTimerQueue` – Creates a timer queue, less common for escalation.

### 5. **Registry & Security Policies**

- `advapi32 LsaOpenPolicy` – Opens a local security policy, useful for attacks requiring system-level security management.
- `advapi32 RegistryCreateKeyEx` – Creates a new registry key, used for manipulating configurations.
- `advapi32 RegistrySetValue` – Sets a value in the registry key, useful for accessing system settings.
- `advapi32 RegistryGetSecurity` – Retrieves security info on a registry key.
- `advapi32 RegistrySetSecurity` – Sets security info on a registry key.
- `advapi32 RegistryKeyDeref` – Delegates a registry key.

---
