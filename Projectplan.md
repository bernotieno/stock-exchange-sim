# 🛠️ 2-Week Work Plan – Stock-exchange-simulator

## ⏲️ Assumptions
- 1 day = 4 working hours
- Timeline: 10 working days (2 weeks)

---

## 📅 Schedule Breakdown

### 🗓️ Week 1

#### Day 1 – Project Understanding & Setup
- Review the entire project brief thoroughly
- Define data structures for stock, processes, and optimization targets
- Choose language and set up project folder/module structure
- Draft input file format based on specification

#### Day 2 – Build Configuration File Parser
- Implement parser for initial stock (`stock_name:quantity`)
- Implement parser for processes and their structure
- Implement parser for the `optimize:(...)` directive
- Test parsing with sample input and print the loaded data

#### Day 3 – Design & Implement Scheduling Framework
- Implement simulation cycle structure
- Design the structure to track running tasks and timing
- Implement check for available stock before starting a process
- Prepare memory structure to hold scheduled tasks

#### Day 4 – Implement Parallel Scheduling Logic
- Implement scheduling of multiple processes per cycle
- Handle process duration, start time, and finish time tracking
- Update stocks once process completes
- Ensure proper dependency management for tasks

#### Day 5 – Generate Log Output
- Build formatted output of schedule in `<cycle>:<process_name>` format
- Write to `.log` file named after input file
- Print final stock state and ending cycle in readable format
- Test with a variety of input files

---

### 🗓️ Week 2

#### Day 6 – Develop the Checker Program (Part 1)
- Load configuration file inside checker
- Load log file and parse lines step-by-step
- Re-simulate execution cycle-by-cycle using log data

#### Day 7 – Develop the Checker Program (Part 2)
- At each cycle, validate stock sufficiency before process execution
- Update stock and check for errors like:
  - Insufficient resources
  - Unknown process
- Display validation success or the specific failure point
- Test with valid and invalid logs

#### Day 8 – Create Custom Configuration Files
- Create File 1: Finite system with limited stock
- Create File 2: Self-regenerating system with infinite loop potential
- Run both through scheduler and checker to ensure correctness
- Polish formatting and organization of files

#### Day 9 – Optimization Strategy Implementation
- Prioritize tasks based on optimization target (e.g. time, product)
- Implement rule to order process execution if multiple are possible
- Add optional support for `serial` and `parallel` strategies
- Test and compare outcomes under different rules

#### Day 10 – Documentation & Final Testing
- Write full usage guide for both programs (scheduler and checker)
- Document input/output formats
- Add comments and cleanup code structure
- Run integration tests across all parts
- Finalize `README.md` with example runs and output

---

## 🧩 Optional Enhancements (if time permits)

- Add command-line flags for verbose output or optimization mode
- Create ASCII visualization of schedule timeline
- Add benchmark comparisons for different strategies
- Handle circular process dependencies more elegantly

---

## 🔚 Final Deliverables

- Stock Exchange Scheduler (`stock_exchange`)
- Checker Program (`checker`)
- Two valid configuration files (`finite.conf`, `looping.conf`)
- Matching `.log` files
- Well-documented `README.md`

