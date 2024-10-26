import http from "k6/http";
import { check, sleep, group } from "k6";
import grpc from "k6/net/grpc";
import { Rate, Counter } from "k6/metrics";

// Define custom metrics
const httpErrors = new Rate("http_errors");
const grpcErrors = new Rate("grpc_req_errors");
const totalGrpcRequests = new Counter("grpc_reqs");

// Initialize gRPC client in the init context
const client = new grpc.Client();
client.load(["../../proto"], "bee_management.proto"); // Adjust path based on your project structure

// Pre-created entity storage
const preCreatedEntities = {
  regions: [],
  apiaries: [],
  hives: [],
  beeCommunities: [],
};

export const options = {
  stages: [
    { duration: "1s", target: 50 },
    { duration: "1s", target: 50 },
    { duration: "1s", target: 1000 },
    { duration: "1s", target: 100 },
    { duration: "1s", target: 20 },
  ],
  thresholds: {
    http_req_failed: ["rate<0.01"], // HTTP request error rate below 1%
    http_req_duration: ["p(95)<600"], // 95% of HTTP requests should be below 500ms
    grpc_req_errors: ["rate<0.01"], // gRPC error rate below 1%
  },
};

// Pre-create entities
export function setup() {
  const BASE_URL = "http://localhost:4040";
  const headers = { "Content-Type": "application/json" };

  console.log(
    "Starting setup: Creating Regions, Apiaries, Hives, and Bee Communities"
  );

  // 1. Create Regions
  for (let i = 1; i <= 5; i++) {
    const res = http.post(
      `${BASE_URL}/region`,
      JSON.stringify({
        name: `Region_${i}`,
        climate_zone: `Climate_${i}`,
      }),
      { headers }
    );

    const success = check(res, {
      [`Create Region_${i} status is 200 or 201`]: (r) =>
        r.status === 200 || r.status === 201,
    });

    if (success) {
      const body = JSON.parse(res.body);
      preCreatedEntities.regions.push(body.region_id || body.id);
      console.log(
        `Created Region_${i} with ID: ${preCreatedEntities.regions[i - 1]}`
      );
    } else {
      console.error(`Failed to create Region_${i}`);
    }
  }

  // 2. Create Apiaries
  for (let i = 1; i <= 10; i++) {
    const region_id =
      preCreatedEntities.regions[i % preCreatedEntities.regions.length];
    const res = http.post(
      `${BASE_URL}/apiary`,
      JSON.stringify({
        location: `Apiary_Location_${i}`,
        manager_id: 1,
        establishment_date: "2023-01-01T15:04:05Z",
      }),
      { headers }
    );

    const success = check(res, {
      [`Create Apiary_${i} status is 200 or 201`]: (r) =>
        r.status === 200 || r.status === 201,
    });

    if (success) {
      const body = JSON.parse(res.body);
      preCreatedEntities.apiaries.push(body.apiary_id || body.id);
      console.log(
        `Created Apiary_${i} with ID: ${preCreatedEntities.apiaries[i - 1]}`
      );
    } else {
      console.error(`Failed to create Apiary_${i}`);
    }
  }

  // 3. Create Hives
  for (let i = 1; i <= 20; i++) {
    const apiary_id =
      preCreatedEntities.apiaries[i % preCreatedEntities.apiaries.length];
    const res = http.post(
      `${BASE_URL}/hive`,
      JSON.stringify({
        apiary_id: apiary_id,
        hive_type: `Type_${i}`,
        installation_date: "2023-01-01T15:04:05Z",
        current_status: "Active",
      }),
      { headers }
    );

    const success = check(res, {
      [`Create Hive_${i} status is 200 or 201`]: (r) =>
        r.status === 200 || r.status === 201,
    });

    if (success) {
      const body = JSON.parse(res.body);
      preCreatedEntities.hives.push(body.hive_id || body.id);
      console.log(
        `Created Hive_${i} with ID: ${preCreatedEntities.hives[i - 1]}`
      );
    } else {
      console.error(`Failed to create Hive_${i}`);
    }
  }

  // 4. Create Bee Communities
  for (let i = 1; i <= 20; i++) {
    const hive_id =
      preCreatedEntities.hives[i % preCreatedEntities.hives.length];
    const res = http.post(
      `${BASE_URL}/bee-community`,
      JSON.stringify({
        hive_id: hive_id,
        queen_age: 2,
        population_estimate: 100,
        health_status: "Healthy",
      }),
      { headers }
    );

    const success = check(res, {
      [`Create BeeCommunity_${i} status is 200 or 201`]: (r) =>
        r.status === 200 || r.status === 201,
    });

    if (success) {
      const body = JSON.parse(res.body);
      preCreatedEntities.beeCommunities.push(body.community_id || body.id);
      console.log(
        `Created BeeCommunity_${i} with ID: ${
          preCreatedEntities.beeCommunities[i - 1]
        }`
      );
    } else {
      console.error(`Failed to create BeeCommunity_${i}`);
    }
  }

  console.log("Setup completed successfully.");
  return preCreatedEntities;
}

// Main load test function
export default function (data) {
  const BASE_URL = "http://localhost:4040";
  const headers = { "Content-Type": "application/json" };

  // Assign entities based on VU number
  const vu = __VU;
  const iter = __ITER;

  // Determine a unique index for the VU
  const regionIndex = (vu - 1) % data.regions.length;
  const apiaryIndex = (vu - 1) % data.apiaries.length;
  const hiveIndex = (vu - 1) % data.hives.length;
  const beeCommIndex = (vu - 1) % data.beeCommunities.length;

  const region_id = data.regions[regionIndex];
  const apiary_id = data.apiaries[apiaryIndex];
  const hive_id = data.hives[hiveIndex];
  const bee_comm_id = data.beeCommunities[beeCommIndex];

  group("REST CRUD Operations", () => {
    // 1. Create Observation Log
    let createObsRes = http.post(
      `${BASE_URL}/observation-log`,
      JSON.stringify({
        hive_id: hive_id,
        observation_date: "2023-07-01T10:00:00Z",
        description: "Initial observation.",
        recommendations: "Monitor closely.",
      }),
      { headers }
    );

    const createObsCheck = check(createObsRes, {
      "Create ObservationLog status is 200 or 201": (r) =>
        r.status === 200 || r.status === 201,
    });

    if (!createObsCheck) {
      httpErrors.add(1);
      return;
    }

    let observationLogId;
    try {
      const respBody = JSON.parse(createObsRes.body);
      observationLogId = respBody.log_id || respBody.id;
    } catch (e) {
      console.error(`Failed to parse Create ObservationLog response: ${e}`);
      httpErrors.add(1);
      return;
    }

    // 2. Get Observation Log
    let getObsRes = http.get(`${BASE_URL}/observation-log/${observationLogId}`);
    check(getObsRes, {
      "Get ObservationLog status is 200": (r) => r.status === 200,
    }) || httpErrors.add(1);

    // 3. Update Observation Log
    let updateObsRes = http.put(
      `${BASE_URL}/observation-log`,
      JSON.stringify({
        log_id: observationLogId,
        hive_id: hive_id,
        observation_date: "2023-07-01T10:00:00Z",
        description: "Updated observation.",
        recommendations: "Continue monitoring.",
      }),
      { headers }
    );

    check(updateObsRes, {
      "Update ObservationLog status is 200": (r) => r.status === 200,
    }) || httpErrors.add(1);

    // 4. Delete Observation Log
    let deleteObsRes = http.del(
      `${BASE_URL}/observation-log/${observationLogId}`
    );
    check(deleteObsRes, {
      "Delete ObservationLog status is 200 or 204": (r) =>
        r.status === 200 || r.status === 204,
    }) || httpErrors.add(1);
  });

  group("gRPC CRUD Operations", () => {
    client.connect("localhost:50051", { plaintext: true });

    try {
      // 1. GetTotalHoneyHarvested
      totalGrpcRequests.add(1);
      const honeyResponse = client.invoke(
        "bee_management.BeeManagementService/GetTotalHoneyHarvested",
        {
          hive_id: hive_id,
          start_date: "2023-01-01",
          end_date: "2024-12-31",
        }
      );

      check(honeyResponse, {
        "GetTotalHoneyHarvested status is OK": (r) =>
          r && r.status === grpc.StatusOK,
      }) || grpcErrors.add(1);

      // 2. AddObservation
      totalGrpcRequests.add(1);
      const addObsResponse = client.invoke(
        "bee_management.BeeManagementService/AddObservation",
        {
          hive_id: hive_id,
          observation_date: "2023-07-01",
          description: "Healthy hive.",
          recommendations: "Maintain current practices.",
        }
      );

      check(addObsResponse, {
        "AddObservation status is OK": (r) => r && r.status === grpc.StatusOK,
      }) || grpcErrors.add(1);

      // 3. GetCommunityHealthStatus
      totalGrpcRequests.add(1);
      const healthResponse = client.invoke(
        "bee_management.BeeManagementService/GetCommunityHealthStatus",
        {
          community_id: bee_comm_id,
        }
      );

      check(healthResponse, {
        "GetCommunityHealthStatus status is OK": (r) =>
          r && r.status === grpc.StatusOK,
      }) || grpcErrors.add(1);

      // 4. UpdateHiveStatus
      totalGrpcRequests.add(1);
      const statusResponse = client.invoke(
        "bee_management.BeeManagementService/UpdateHiveStatus",
        {
          hive_id: hive_id,
          new_status: "Active",
        }
      );

      check(statusResponse, {
        "UpdateHiveStatus status is OK": (r) => r && r.status === grpc.StatusOK,
      }) || grpcErrors.add(1);
    } catch (err) {
      console.error(`gRPC error: ${err}`);
      grpcErrors.add(1);
    } finally {
      client.close();
    }
  });

  // Pause between iterations
  sleep(Math.random() + 1);
}
