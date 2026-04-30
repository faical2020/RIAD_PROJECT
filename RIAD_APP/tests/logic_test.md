The tests were grouped in TestHybridLogic and consisted of three main parts:
1.  RoomValidation:
    *   Tested the business logic in pkg/logic/room.go.
    *   Verified that valid room data (positive number, non-negative price) passes.
    *   Verified that invalid data (e.g., negative room number) is correctly rejected.
    *   Goal: Ensure consistent validation regardless of whether the app is offline or online.
2.  LocalPersistence:
    *   Tested the internal/db SQLite implementation.
    *   Saved a room to the local database.
    *   Retrieved the list of rooms and checked if the saved room existed.
    *   Goal: Confirm that the "Local" part of "Local-First" actually works (data survives in the DB).
3.  SyncFlagging:
    *   Tested the synchronization state management.
    *   Saved a reservation and verified it was marked as synced = 0 (unsynced) by default.
    *   Called MarkSynced() and verified that the record was no longer returned by GetUnsynced().
    *   Goal: Ensure the app can correctly identify which records need to be pushed to the server when connectivity is restored.
I should explain this concisely.
I performed a suite of automated tests in tests/logic_test.go to verify the "Local-First" engine without needing to launch the full GUI. 
Here is the breakdown of what was tested:
1. Business Logic Validation (RoomValidation)
I tested the rules defined in pkg/logic. This ensures that the app rejects invalid data before it even reaches the database or the server.
*   Test: Tried saving a room with a negative room number and a negative price.
*   Result: The logic correctly blocked these entries, preventing corrupted data from entering the system.
2. Local Persistence (LocalPersistence)
I verified that the SQLite integration is working correctly.
*   Test: I wrote a room record to the local test_riad.db and then immediately tried to read it back.
*   Result: The data was retrieved exactly as saved, confirming that the app can store and read data while completely offline.
3. Synchronization State (SyncFlagging)
This is the most critical part of the hybrid architecture. I tested the "Pending Sync" mechanism.
*   Test: 
    1.  Created a reservation $\rightarrow$ Verified it was automatically marked as Unsynced (synced = 0).
    2.  Simulated a successful API push $\rightarrow$ Marked the record as Synced (synced = 1).
    3.  Checked the "Unsynced" list again $\rightarrow$ Verified the record had disappeared from the pending list.
*   Result: Confirmed that the app can accurately track which local changes still need to be pushed to the RIAD_SERVER.