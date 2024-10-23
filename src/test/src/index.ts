import * as grpc from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";
import path from "node:path";

// Define the path to the proto file
const PROTO_PATH = path.join(__dirname, "../../proto/bee_management.proto");

// Load the protobuf
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

// Load the package definition
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;

// Get the BeeManagementService
const beeManagement = protoDescriptor.bee_management.BeeManagementService;

// Create a client instance
const client = new beeManagement("localhost:50051", grpc.credentials.createInsecure());

// Helper function to promisify client methods
function promisifyClientMethod(method: Function) {
  return (...args: any[]) => {
    return new Promise((resolve, reject) => {
      method(...args, (error: any, response: any) => {
        if (error) {
          reject(error);
        } else {
          resolve(response);
        }
      });
    });
  };
}

// Promisified client methods
const getTotalHoneyHarvested = promisifyClientMethod(client.GetTotalHoneyHarvested.bind(client));
const addObservation = promisifyClientMethod(client.AddObservation.bind(client));
const getCommunityHealthStatus = promisifyClientMethod(
  client.GetCommunityHealthStatus.bind(client),
);
const updateHiveStatus = promisifyClientMethod(client.UpdateHiveStatus.bind(client));
const getAvgTemperature = promisifyClientMethod(client.GetAvgTemperature.bind(client));
const assignMaintenancePlan = promisifyClientMethod(client.AssignMaintenancePlan.bind(client));
const hasRegionAccess = promisifyClientMethod(client.HasRegionAccess.bind(client));
const registerIncident = promisifyClientMethod(client.RegisterIncident.bind(client));
const getLatestSensorReading = promisifyClientMethod(client.GetLatestSensorReading.bind(client));
const createProductionReport = promisifyClientMethod(client.CreateProductionReport.bind(client));

async function main() {
  try {
    // 1. Get Total Honey Harvested
    const totalHoney = await getTotalHoneyHarvested({
      hive_id: 1,
      start_date: "2023-01-01",
      end_date: "2024-12-31",
    });
    console.log("Total Honey Harvested:", totalHoney.total_honey);

    // 2. Add Observation
    const addObsResponse = await addObservation({
      hive_id: 1,
      observation_date: "2023-04-15",
      description: "Queen is healthy.",
      recommendations: "Continue current beekeeping practices.",
    });
    console.log("Add Observation Response:", addObsResponse);

    // 3. Get Community Health Status
    const communityHealth = await getCommunityHealthStatus({
      community_id: 1,
    });
    console.log("Community Health Status:", communityHealth.health_status);

    // 4. Update Hive Status
    const updateStatusResponse = await updateHiveStatus({
      hive_id: 1,
      new_status: "Active",
    });
    console.log("Update Hive Status Response:", updateStatusResponse);

    // 5. Get Average Temperature
    const avgTemp = await getAvgTemperature({
      region_id: 5,
      days: 30,
    });
    console.log("Average Temperature:", avgTemp.avg_temperature);

    // 6. Assign Maintenance Plan
    const assignPlanResponse = await assignMaintenancePlan({
      plan_id: 7,
      user_id: 3,
    });
    console.log("Assign Maintenance Plan Response:", assignPlanResponse);

    // 7. Has Region Access
    const regionAccess = await hasRegionAccess({
      user_id: 42,
      region_id: 5,
    });
    console.log("Has Region Access:", regionAccess.has_access);

    // 8. Register Incident
    const registerIncidentResponse = await registerIncident({
      hive_id: 1,
      incident_date: "2023-05-20",
      description: "Varroa mite infestation detected.",
      severity: "High",
    });
    console.log("Register Incident Response:", registerIncidentResponse);

    // 9. Get Latest Sensor Reading
    const latestSensor = await getLatestSensorReading({
      hive_id: 1,
      sensor_type: "humidity",
    });
    console.log("Latest Sensor Reading:", latestSensor);

    // 10. Create Production Report
    const createReportResponse = await createProductionReport({
      apiary_id: 1,
      start_date: "2023-01-01",
      end_date: "2023-06-30",
    });
    console.log("Create Production Report Response:", createReportResponse);
  } catch (error) {
    console.error("An error occurred:", error);
  } finally {
    client.close();
  }
}

main();
