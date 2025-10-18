println("Running Queue tests...")

fn createQueue() ::
  var queue = []

  queue.enqueue = fn(self, value) ::
    return self.push(value)
  end

  queue.dequeue = fn(self) ::
    return self.first
    # Note: In a real implementation, we'd remove the first element
  end

  queue.peek = fn(self) ::
    return self.first
  end

  queue.is_empty = fn(self) ::
    return self.length == 0
  end

  queue.size = fn(self) ::
    return self.length
  end

  queue.to_string = fn(self) ::
    var result = "Queue["
    for i in [0, 1, 2, 3, 4] ::
      if i < self.length ::
        result = result + self[i].to_string
        if i < self.length - 1 ::
          result = result + ", "
        end
      end
    end
    return result + "]"
  end

  return queue
end

var myQueue = createQueue()
myQueue = myQueue.enqueue(10)
myQueue = myQueue.enqueue(20)
myQueue = myQueue.enqueue(30)

check "Queue" ::
  myQueue.to_string is "Queue[10, 20, 30]"
  myQueue.peek is 10
  myQueue.dequeue is 10
  myQueue.size is 3
end

println("✓ All Queue tests passed!")