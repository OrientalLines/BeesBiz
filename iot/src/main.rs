use actix::{Actor, Context, Handler, Message};

// Define the state of your "GenServer"
struct Counter {
    count: i32,
}

// Implement Actor trait for your state
impl Actor for Counter {
    type Context = Context<Self>;

    fn started(&mut self, _ctx: &mut Self::Context) {
        println!("Counter actor started");
    }
}

// Define messages that your actor can handle
#[derive(Message)]
#[rtype(result = "i32")]
struct Increment;

#[derive(Message)]
#[rtype(result = "i32")]
struct GetCount;

// Implement handlers for your messages
impl Handler<Increment> for Counter {
    type Result = i32;

    fn handle(&mut self, _msg: Increment, _ctx: &mut Context<Self>) -> Self::Result {
        self.count += 1;
        self.count
    }
}

impl Handler<GetCount> for Counter {
    type Result = i32;

    fn handle(&mut self, _msg: GetCount, _ctx: &mut Context<Self>) -> Self::Result {
        self.count
    }
}

// Usage example
#[actix::main]
async fn main() {
    // Start the actor
    let addr = Counter { count: 0 }.start();

    // Send messages to the actor
    let result = addr.send(Increment).await.unwrap();
    println!("After increment: {}", result);

    let count = addr.send(GetCount).await.unwrap();
    println!("Current count: {}", count);
}
