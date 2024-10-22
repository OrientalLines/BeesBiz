import http from "k6/http";
import { check, sleep } from "k6";
import { Counter } from "k6/metrics";

export const errorCount = new Counter("errors");

const BASE_URL = "http://localhost:4040"; // Replace with your server's base URL

export default function () {
  // Apiary
  let apiaryPayload = {
    location: "Test Location",
    manager_id: 1,
    establishment_date: "2023-01-01",
  };

  // Create Apiary
  let res = http.post(`${BASE_URL}/apiary`, JSON.stringify(apiaryPayload), {
    headers: { "Content-Type": "application/json" },
  });

  let createdApiary = res.json();
  check(res, {
    "Create Apiary status is 200": (r) => r.status === 200,
    "Create Apiary has id": (r) => r.json("id") !== undefined,
  }) || errorCount.add(1);

  let apiaryId = createdApiary.id;

  // Get Apiary
  res = http.get(`${BASE_URL}/apiary/${apiaryId}`);
  check(res, {
    "Get Apiary status is 200": (r) => r.status === 200,
    "Get Apiary has correct id": (r) => r.json("id") === apiaryId,
  }) || errorCount.add(1);

  // Update Apiary
  apiaryPayload.location = "Updated Location";
  res = http.put(`${BASE_URL}/apiary`, JSON.stringify(apiaryPayload), {
    headers: { "Content-Type": "application/json" },
  });
  check(res, {
    "Update Apiary status is 200": (r) => r.status === 200,
    "Update Apiary location updated": (r) =>
      r.json("location") === "Updated Location",
  }) || errorCount.add(1);

  // Get All Apiaries
  res = http.get(`${BASE_URL}/apiaries`);
  check(res, {
    "Get All Apiaries status is 200": (r) => r.status === 200,
    "Get All Apiaries is array": (r) => Array.isArray(r.json()),
  }) || errorCount.add(1);

  // Delete Apiary
  res = http.del(`${BASE_URL}/apiary/${apiaryId}`);
  check(res, {
    "Delete Apiary status is 204": (r) => r.status === 204,
  }) || errorCount.add(1);

  // Repeat similar steps for other entities...

  // Example for User
  let userPayload = {
    name: "Test User",
    email: "testuser@example.com",
    role: "admin",
  };

  // Create User
  res = http.post(`${BASE_URL}/user`, JSON.stringify(userPayload), {
    headers: { "Content-Type": "application/json" },
  });

  let createdUser = res.json();
  check(res, {
    "Create User status is 200": (r) => r.status === 200,
    "Create User has id": (r) => r.json("id") !== undefined,
  }) || errorCount.add(1);

  let userId = createdUser.id;

  // Get User
  res = http.get(`${BASE_URL}/user/${userId}`);
  check(res, {
    "Get User status is 200": (r) => r.status === 200,
    "Get User has correct id": (r) => r.json("id") === userId,
  }) || errorCount.add(1);

  // Update User
  userPayload.name = "Updated User";
  res = http.put(`${BASE_URL}/user`, JSON.stringify(userPayload), {
    headers: { "Content-Type": "application/json" },
  });
  check(res, {
    "Update User status is 200": (r) => r.status === 200,
    "Update User name updated": (r) => r.json("name") === "Updated User",
  }) || errorCount.add(1);

  // Get All Users
  res = http.get(`${BASE_URL}/users`);
  check(res, {
    "Get All Users status is 200": (r) => r.status === 200,
    "Get All Users is array": (r) => Array.isArray(r.json()),
  }) || errorCount.add(1);

  // Delete User
  res = http.del(`${BASE_URL}/user/${userId}`);
  check(res, {
    "Delete User status is 204": (r) => r.status === 204,
  }) || errorCount.add(1);

  // Add similar blocks for other entities like Hive, Incident, Sensor, etc.

  sleep(1); // Pause between iterations
}
