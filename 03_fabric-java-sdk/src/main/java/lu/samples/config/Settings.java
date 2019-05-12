package lu.samples.config;

import static lu.samples.config.ConfFactory.*;
import java.util.Optional;
import lu.samples.config.ConfFactory.Conf;

/**
 * @TODO
 *
 * @author Emilien Charton
 * @date 14 nov. 2018
 */
public enum Settings {
	
	ORG1_USER_ADMIN("admin.warehouse", conf()
									.prop("domain", "Admin@warehouse.sfeir.lu")
									.prop("name", "Admin@warehouse.sfeir.lu")
									.prop("pwd", "adminpw")
									.prop("mspId", "WarehouseMSP")
									.instance()),
	
	ORG1_PEER0("peer0.warehouse", conf()
								.prop("name", "peer0.warehouse.sfeir.lu")
								.prop("domain", "warehouse.sfeir.lu")
								.prop("url", "grpc://0.0.0.0:8051")
								.instance()),
	
	ORG1_CA("ca.warehouse", conf()
							.prop("name", "ca.warehouse.sfeir.lu")
							.prop("url", "http://0.0.0.0:8054")
							.prop("domain", "warehouse.sfeir.lu")
							.instance()),
	
	
	ORDERER("orderer", conf()
							.prop("name", "orderer.sfeir.lu")
							.prop("url", "grpc://0.0.0.0:7050")
							.instance()),
	
	CHANNEL("channel", conf()
							.prop("name", "sfeircn")
							.instance())
	
	;
	/**
	 * Reference on the name of setting purpose.
	 */
	private String name;
	/**
	 * Reference on the <code>Conf</code> of setting purpose.
	 */
	private Conf configuration;

	/**
	 * Default constructor.
	 */
	Settings(final String name, final Conf configuration) {
		this.name=name;
		this.configuration=configuration;
	}
	
	
	/**
	 * Accessor in reading on the name.
	 * @return the name.
	 */
	public String getName() {
		return name;
	}

	/**
	 * Accessor in reading on the configuration.
	 * @return the configuration.
	 */
	public Conf getConfiguration() {
		return configuration;
	}
	
	/**
	 * Accessor in reading on the property value.
	 * @param k The key.
	 * @return The optional value.
	 */
	public Optional<String> val(final String k) {
		
		return (configuration!=null) ? configuration.val(k) : Optional.empty();
	}
	
}
