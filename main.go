package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/sql"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		project := config.Get(ctx, "gcp:project")
		if project == "" {
			return fmt.Errorf("gcp:project is required")
		}

		region := config.Get(ctx, "gcp:region")
		if region == "" {
			region = "us-central1"
		}

		instanceName := config.Get(ctx, "cloudSql:instanceName")
		if instanceName == "" {
			return fmt.Errorf("cloudSql:instanceName is required")
		}

		databaseName := config.Get(ctx, "cloudSql:databaseName")
		if databaseName == "" {
			return fmt.Errorf("cloudSql:databaseName is required")
		}

		databaseUser := config.Get(ctx, "cloudSql:databaseUser")
		if databaseUser == "" {
			return fmt.Errorf("cloudSql:databaseUser is required")
		}

		databaseVersion := config.Get(ctx, "cloudSql:databaseVersion")
		if databaseVersion == "" {
			databaseVersion = "MYSQL_8_0"
		}

		tier := config.Get(ctx, "cloudSql:tier")
		if tier == "" {
			tier = "db-f1-micro"
		}

		// Generate random password
		dbPassword, err := random.NewRandomPassword(ctx, "db-password", &random.RandomPasswordArgs{
			Length:          pulumi.Int(16),
			Special:         pulumi.Bool(true),
			OverrideSpecial: pulumi.String("!@#$%^&*"),
		})
		if err != nil {
			return err
		}

		// Create Cloud SQL instance
		instance, err := sql.NewDatabaseInstance(ctx, "instance", &sql.DatabaseInstanceArgs{
			Name:               pulumi.String(instanceName),
			DatabaseVersion:    pulumi.String(databaseVersion),
			Region:             pulumi.String(region),
			Project:            pulumi.String(project),
			DeletionProtection: pulumi.Bool(false),
			Settings: &sql.DatabaseInstanceSettingsArgs{
				Tier: pulumi.String(tier),
				IpConfiguration: &sql.DatabaseInstanceSettingsIpConfigurationArgs{
					SslMode: pulumi.String("ALLOW_UNENCRYPTED_AND_ENCRYPTED"),
					AuthorizedNetworks: sql.DatabaseInstanceSettingsIpConfigurationAuthorizedNetworkArray{
						&sql.DatabaseInstanceSettingsIpConfigurationAuthorizedNetworkArgs{
							Value: pulumi.String("0.0.0.0/0"),
						},
					},
				},
				BackupConfiguration: &sql.DatabaseInstanceSettingsBackupConfigurationArgs{
					Enabled: pulumi.Bool(true),
				},
				AvailabilityType: pulumi.String("ZONAL"),
			},
		})
		if err != nil {
			return err
		}

		// Create database
		database, err := sql.NewDatabase(ctx, "database", &sql.DatabaseArgs{
			Name:     pulumi.String(databaseName),
			Instance: instance.Name,
			Project:  pulumi.String(project),
		})
		if err != nil {
			return err
		}

		// Create user
		user, err := sql.NewUser(ctx, "user", &sql.UserArgs{
			Name:     pulumi.String(databaseUser),
			Instance: instance.Name,
			Password: dbPassword.Result,
			Project:  pulumi.String(project),
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("instanceName", instance.Name)
		ctx.Export("connectionName", instance.ConnectionName)
		ctx.Export("publicIpAddress", instance.PublicIpAddress)
		ctx.Export("databaseName", database.Name)
		ctx.Export("databaseUser", user.Name)
		ctx.Export("databasePassword", dbPassword.Result)

		return nil
	})
}
