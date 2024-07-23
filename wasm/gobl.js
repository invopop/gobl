var ready = false;
var queue = [];
var inFlight = {};
var req_id = 0;

const worker = new window.Worker("worker.js");

worker.onmessage = (event) => {
  if (event.data && event.data.ready) {
    console.log("worker is ready");
    ready = true;
    for (var i = 0; i < queue.length; i++) {
      worker.postMessage(queue[i]);
    }
    return true;
  }
  const waiting = inFlight[event.data.req_id];
  delete inFlight[event.data.req_id];
  if (!waiting) {
    console.log(
      "got a response for an unregistered request: " + event.data.req_id
    );
    return true;
  }
  if (event.data.error) {
    waiting.reject(event.data.error);
    return true;
  }
  waiting.resolve(event.data.payload);
};

function sendMessage(data) {
  if (!data.req_id) {
    data.req_id = `req${++req_id}`;
  }
  const promise = new Promise((resolve, reject) => {
    inFlight[data.req_id] = {
      resolve: resolve,
      reject: reject,
    };
    // resolve("foo");
  });
  if (!ready) {
    queue.push(data);
    return promise;
  }
  worker.postMessage(data);
  return promise;
}

const keygen = async function (opts) {
  const indent = delete opts.indent;
  return sendMessage({
    action: "keygen",
    indent: indent,
  });
};

const build = async function (opts) {
  const indent = delete opts.indent;
  return sendMessage({
    action: "build",
    payload: opts,
    indent: indent,
  });
};

const verify = async function (opts) {
  const indent = delete opts.indent;
  return sendMessage({
    action: "verify",
    payload: opts,
    indent: indent,
  });
};

export { keygen, build, verify };
