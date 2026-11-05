"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getStatistics = getStatistics;
const matrices_validator_1 = require("../validators/matrices.validator");
const statistics_service_1 = require("../services/statistics.service");
function getStatistics(req, res) {
    const matrices = (0, matrices_validator_1.validateBody)(req.body);
    res.json((0, statistics_service_1.computeStatistics)(matrices));
}
