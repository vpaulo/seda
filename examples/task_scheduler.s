## Task Scheduler - Real-World Time Module Example
## Demonstrates: Time manipulation, arrays, maps, functions, control flow

## Task structure: {name, created, dueDate, status, priority}
var tasks = []

## Create a new task
fn create_task(name, due_days_from_now, priority) ::
    var now = Time.now()
    var due_date = now.add_days(due_days_from_now)

    var task = {
        "name": name,
        "created": now,
        "dueDate": due_date,
        "status": "pending",
        "priority": priority
    }

    var _ = tasks.push(task)
    println("Created task: #{name} (due in #{due_days_from_now} days)")
    return task
end

## Check if a task is overdue
fn is_overdue(task) ::
    var now = Time.now()
    return now.is_after(task.dueDate)
end

## Get days until due date
fn days_until_due(task) ::
    var now = Time.now()
    var diff_seconds = task.dueDate.diff(now)
    var diff_days = diff_seconds / 86400  ## Convert seconds to days
    return diff_days.floor()
end

## Complete a task
fn complete_task(task_name) ::
    for task in tasks ::
        if task.name == task_name ::
            task.status = "completed"
            println("Completed task: #{task_name}")
            return task
        end
    end
    println("Task not found: #{task_name}")
    return nil
end

## List all tasks with status
fn list_tasks() ::
    println("\n==================================================")
    println("Task List - #{tasks.length()} tasks")
    println("==================================================\n")

    if tasks.length() == 0 ::
        println("No tasks yet!")
        return nil
    end

    for task in tasks ::
        var status_icon = case task.status ::
            "completed" => "✓"
            "pending" => "○"
            _ => "?"
        end

        var priority_str = case task.priority ::
            1 => "[HIGH]"
            2 => "[MED]"
            3 => "[LOW]"
            _ => "[?]"
        end

        var days_left = days_until_due(task)
        var overdue = is_overdue(task)

        var due_str = ""
        if overdue ::
            due_str = "OVERDUE by #{-days_left} days"
        else ::
            due_str = "#{days_left} days left"
        end

        println("#{status_icon} #{priority_str} #{task.name}")
        var created_str = task.created.format("YYYY-MM-DD HH:mm")
        var due_date_str = task.dueDate.format("YYYY-MM-DD")
        println("   Created: #{created_str}")
        println("   Due: #{due_date_str} (#{due_str})")
        println("")
    end
end

## Get statistics about tasks
fn get_statistics() ::
    var total = tasks.length()
    var completed = 0
    var pending = 0
    var overdue = 0
    var high_priority = 0

    for task in tasks ::
        if task.status == "completed" ::
            completed = completed + 1
        else ::
            pending = pending + 1
            if is_overdue(task) ::
                overdue = overdue + 1
            end
        end

        if task.priority == 1 ::
            high_priority = high_priority + 1
        end
    end

    return {
        "total": total,
        "completed": completed,
        "pending": pending,
        "overdue": overdue,
        "highPriority": high_priority
    }
end

## Print statistics
fn print_statistics() ::
    var stats = get_statistics()

    println("\n==================================================")
    println("Task Statistics")
    println("==================================================")
    println("Total tasks:      #{stats.total}")
    println("Completed:        #{stats.completed}")
    println("Pending:          #{stats.pending}")
    println("Overdue:          #{stats.overdue}")
    println("High priority:    #{stats.highPriority}")

    if stats.total > 0 ::
        var completion_pct = (stats.completed * 100) / stats.total
        var completion_rate = completion_pct.floor()
        println("Completion rate:  #{completion_rate}%")
    end
    println("")
end

## Get overdue tasks
fn get_overdue_tasks() ::
    var overdue_tasks = []

    for task in tasks ::
        if task.status != "completed" && is_overdue(task) ::
            var _ = overdue_tasks.push(task)
        end
    end

    return overdue_tasks
end

## Print overdue tasks
fn print_overdue_tasks() ::
    var overdue_tasks = get_overdue_tasks()

    if overdue_tasks.length() == 0 ::
        println("\nNo overdue tasks! 🎉")
        return nil
    end

    println("\n==================================================")
    println("⚠️  Overdue Tasks - #{overdue_tasks.length()} task(s)")
    println("==================================================\n")

    for task in overdue_tasks ::
        var days_left = days_until_due(task)
        var days_overdue = -days_left
        var due_date_str = task.dueDate.format("YYYY-MM-DD")
        println("• #{task.name}")
        println("  Due: #{due_date_str} (#{days_overdue} days ago)")
        println("")
    end
end

## Sort tasks by due date
fn sort_tasks_by_due_date() ::
    ## Simple bubble sort by due date
    for i in 0..(tasks.length() - 1) ::
        for j in (i + 1)..tasks.length() ::
            var task_i = tasks[i]
            var task_j = tasks[j]

            if task_i.dueDate.is_after(task_j.dueDate) ::
                ## Swap
                tasks[i] = task_j
                tasks[j] = task_i
            end
        end
    end

    println("Tasks sorted by due date")
end

## Get tasks due this week
fn get_tasks_due_this_week() ::
    var now = Time.now()
    var week_from_now = now.add_days(7)
    var this_week = []

    for task in tasks ::
        if task.status != "completed" ::
            if task.dueDate.is_after(now) && task.dueDate.is_before(week_from_now) ::
                var _ = this_week.push(task)
            end
        end
    end

    return this_week
end

## Print tasks due this week
fn print_tasks_due_this_week() ::
    var this_week = get_tasks_due_this_week()

    if this_week.length() == 0 ::
        println("\nNo tasks due this week")
        return nil
    end

    println("\n==================================================")
    println("Tasks Due This Week - #{this_week.length()} task(s)")
    println("==================================================\n")

    for task in this_week ::
        var days_left = days_until_due(task)
        println("• #{task.name} (in #{days_left} days)")
    end
    println("")
end

## ========================================
## MAIN PROGRAM - Demo the task scheduler
## ========================================

println("Task Scheduler Demo")
println("====================\n")

## Create some sample tasks
create_task("Write documentation", 7, 1)
create_task("Review pull requests", 3, 1)
create_task("Update dependencies", 14, 2)
create_task("Plan sprint", 5, 2)
create_task("Refactor authentication", 21, 3)

## Create an overdue task (simulate by creating it "in the past")
var overdue_task = {
    "name": "Fix critical bug",
    "created": Time.now().add_days(-10),
    "dueDate": Time.now().add_days(-2),
    "status": "pending",
    "priority": 1
}
var _ = tasks.push(overdue_task)

## List all tasks
list_tasks()

## Print statistics
print_statistics()

## Print overdue tasks
print_overdue_tasks()

## Print tasks due this week
print_tasks_due_this_week()

## Complete some tasks
complete_task("Review pull requests")
complete_task("Fix critical bug")

## Sort and display again
println("\nAfter completing tasks:\n")
sort_tasks_by_due_date()
list_tasks()
print_statistics()
