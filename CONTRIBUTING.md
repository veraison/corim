### Contributing

Thank you for your interest in contributing to this project! We welcome and
appreciate all forms of contribution—bug fixes, new features, documentation
improvements, and tests. Following these guidelines helps ensure a smooth and
efficient review process.

---

### Before You Start

1. **Search Existing Work**  
   Check the repository’s [issues] and [pull requests] to see if someone is
   already working on the same idea or problem.

2. **Discuss Large Changes**  
   For major changes or new features, open an [issue] first to discuss the
   design, scope, and implementation plan with the maintainers.

3. **Read the Code of Conduct**  
   Please review and follow our [Code of Conduct](https://github.com/veraison/corim/blob/main/CODE_OF_CONDUCT.md).
   We expect all contributors to maintain respectful and inclusive behavior.

---

### How to Contribute

Follow the steps below to prepare and submit your contribution:

### 1. Fork & Clone

Start by forking the repository to your GitHub account and cloning it locally:

```bash
git clone https://github.com/<your-username>/corim.git
cd corim
```

### 2. Set Up the Project

Refer to the project setup instructions in the `README.md` (or the
repository's setup documentation) to install dependencies and configure your
local environment.

### 3. Create a Branch

Create a new, descriptively named branch for your work. Use one of the
recommended prefixes:

`feature/` for new features.

`fix/` for bug fixes.

`chore/` for maintenance or build-related tasks.

```bash
git checkout -b feature/your-feature-name
```

### 4. Make Changes

**Keep it Focused:** Where possible, keep your changes focused on a single
concern.

**Adhere to Standards:** Follow the project's coding style and conventions.

**Document Code:** Keep your code well-documented and include clear inline
comments where complex logic is involved.

### 5. Add or Update Tests

Where applicable, add tests that fully cover your changes. New or updated
tests are crucial for preventing regressions and significantly expedite the
review process.

### 6. Run Tests

Before submitting your contribution, ensure the entire test suite passes
successfully using the project's designated test command:

```bash
make test
```

### 7. Commit Your Changes

Write clear, concise, and descriptive commit messages.

Example:

```bash
git add .
git commit -m "Fix: handle null pointer in data loader (fixes #42)"
```

### 8. Push Your Changes

Push your new branch to your forked repository:

```bash
git push origin feature/your-feature-name
```

### 9. Create a Pull Request (PR)

Open a pull request from your forked branch to the main repository's target
branch (usually main or master).

**Detailed Description:** Include a comprehensive description of your changes,
the rationale behind them, and any relevant issue numbers (e.g., Closes #101).

**Checklist:** Consider including a small checklist in your PR description
(e.g., tests added/updated, documentation updated, code style checked).

### 10. Prior to Merge

Ensure all CI tests pass 

Ensure the Pull Request is reviewed by at least two maintainers of the project.

Incorporate all the review comments. In case of conflicting comments, 
please schedule a meeting and coordinate among all reviewers to reach a consensus.

Then only request maintainers to Merge the Pull Request
