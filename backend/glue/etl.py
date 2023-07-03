import sys
from awsglue.transforms import *
from awsglue.utils import getResolvedOptions
from pyspark.context import SparkContext
from pyspark.sql import functions as F
from awsglue.context import GlueContext
from awsglue.job import Job
from awsglue.dynamicframe import DynamicFrame

## @params: [JOB_NAME, s3_bucket, s3_key, table, workers]
args = getResolvedOptions(sys.argv, ['JOB_NAME', 's3_bucket', 's3_key', 'table', 'workers'])

#Create spark context
sc = SparkContext()
glueContext = GlueContext(sc)
spark = glueContext.spark_session
job = Job(glueContext)
job.init(args['JOB_NAME'], args)

# Gets the s3 file source
s3_file_source = f's3://{args["s3_bucket"]}/{args["s3_key"]}'

# Extract s3 file longo dynamic frame
inputGDF = glueContext.create_dynamic_frame_from_options(
    connection_type="s3", 
    format="csv", 
    format_options={"withHeader":True},
    connection_options={"paths": [s3_file_source], 
    "recurse": True}
)

# Apply mapping
inputGDF = ApplyMapping.apply(
    frame = inputGDF, 
    mappings = [
        ("accountNumber", "string", "accountNumber", "string"),
        ("customerId", "string", "customerId", "string"),
        ("creditLimit", "string", "creditLimit", "double"),
        ("availableMoney", "string", "availableMoney", "double"),
        ("transactionDateTime", "string", "transactionDateTime", "string"),
        ("transactionAmount", "string", "transactionAmount", "double"),
        ("merchantName", "string", "merchantName", "string"),
        ("acqCountry", "string", "acqCountry", "string"),
        ("merchantCountryCode", "string", "merchantCountryCode", "string"),
        ("posEntryMode", "string", "posEntryMode", "long"),
        ("posConditionCode", "string", "posConditionCode", "long"),
        ("merchantCategoryCode", "string", "merchantCategoryCode", "string"),
        ("currentExpDate", "string", "currentExpDate", "string"),
        ("accountOpenDate", "string", "accountOpenDate", "string"),
        ("dateOfLastAddressChange", "string", "dateOfLastAddressChange", "string"),
        ("cardCVV", "string", "cardCVV", "long"),
        ("enteredCVV", "string", "enteredCVV", "long"),
        ("cardLast4Digits", "string", "cardLast4Digits", "long"),
        ("transactionType", "string", "transactionType", "string"),
        ("currentBalance", "string", "currentBalance", "double"),
        ("cardPresent", "string", "cardPresent", "string"),
        ("isFraud", "string", "isFraud", "string"),
        ("CountryCode", "string", "CountryCode", "string"),
    ]
)

# Convert dynamic frame to data frame
inputDF = inputGDF.toDF()

# Create a UUID column to be used as a primary key
inputDF = inputDF.withColumn("id", F.monotonically_increasing_id())

# Remove everything past hyphen in CountryCode column
inputDF = inputDF.withColumn("CountryCode", inputDF["CountryCode"].substr(0, 2))

# If isFraud is "FALSE", set to "TRUE"
inputDF = inputDF.withColumn("isFraud", F.when(inputDF["isFraud"] == "FALSE", "TRUE").otherwise(inputDF["isFraud"]))

# Convert transactionDateTime yyyy-mm-ddThh:mm:ss to a timestamp
inputDF = inputDF.withColumn("transactionDateTime", F.unix_timestamp(inputDF["transactionDateTime"], "yyyy-MM-dd'T'HH:mm:ss").cast("timestamp"))

# Convert accountOpenDate, dateOfLastAddressChange mm/dd/yy to a timestamp
inputDF = inputDF.withColumn("accountOpenDate", F.unix_timestamp(inputDF["accountOpenDate"], "M/d/yy").cast("timestamp"))
inputDF = inputDF.withColumn("dateOfLastAddressChange", F.unix_timestamp(inputDF["dateOfLastAddressChange"], "M/d/yy").cast("timestamp"))

# Convert dateframe to dynamic frame
inputGDF = DynamicFrame.fromDF(inputDF, glueContext, "inputGDF")

# Load S3 files into DynamoDB
Datasink1 = glueContext.write_dynamic_frame_from_options(
    frame=inputGDF, 
    connection_type="dynamodb", 
    connection_options={
        "dynamodb.output.tableName": args["table"], 
        "dynamodb.throughput.write.percent": "1.0",
        "dynamodb.output.numParallelTasks": args["workers"]
    },
)