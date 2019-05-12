package lu.samples.model;

import java.io.Serializable;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;

import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.User;

/**
 * @TODO
 *
 * @author Emilien Charton
 * @date 14 nov. 2018
 */
public class UserDefault implements Serializable, User {

	/**
	 * Reference on the serialVersionUID.
	 */
	private static final long serialVersionUID = -1310053944589380387L;
	
	/**
	 * Reference on the user name.
	 */
	private String name;
	
	/**
	 * Reference on the user's account.
	 */
	private String account;
	
	/**
	 * Reference on the user roles.
	 */
	private Set<String> roles = new HashSet<>();
	
	/**
	 * Reference on the <code>Enrollment</code>.
	 */
	private Enrollment enrollment;
	
	/**
	 * Reference on the domain.
	 */
	private String domain;
	
	/**
	 * Reference on the msp identifier.
	 */
	private String mspId;
	
	/**
	 * Default constructor.
	 */
	public UserDefault() {
		// Empty.
	}

	/**
	 * Default constructor.
	 */
	public UserDefault(final String name, final String domain, final String msp, final Enrollment enrollment) {
		
		this.name = name;
		this.enrollment = enrollment;
		this.domain = domain;
		this.mspId = msp;
	}
	
	/**
	 * Accessor in reading on the name.
	 * @return the name.
	 */
	public String getName() {
		return name;
	}

	/**
	 * Accessor in writing on the name.
	 * @param name the name.
	 */
	public void setName(final String name) {
		this.name = name;
	}

	/**
	 * Accessor in reading on the mspIdentifier.
	 * @return the mspIdentifier.
	 */
	@Override
	public String getMspId() {
		return mspId;
	}

	/**
	 * Accessor in writing on the mspIdentifier.
	 * @param identifier the mspIdentifier.
	 */
	public void setMspId(final String identifier) {
		this.mspId = identifier;
	}

	/**
	 * Accessor in reading on the roles.
	 * @return the roles.
	 */
	public Set<String> getRoles() {
		return roles;
	}

	/**
	 * Accessor in writing on the roles.
	 * @param roles the roles.
	 */
	public void addRoles(final String ...roleNames) {
		if (roles==null || roleNames.length==0) {
			return;
		}
		Arrays.asList(roleNames).stream().forEach(roles::add);
	}

	/**
	 * Accessor in reading on the enrollment.
	 * @return the enrollment.
	 */
	public Enrollment getEnrollment() {
		return enrollment;
	}

	/**
	 * Accessor in writing on the enrollment.
	 * @param enrollment the enrollment.
	 */
	public void setEnrollment(Enrollment enrollment) {
		this.enrollment = enrollment;
	}

	/**
	 * Accessor in reading on the domain.
	 * @return the domain.
	 */
	@Override
	public String getAffiliation() {
		return domain;
	}

	/**
	 * Accessor in writing on the domain.
	 * @param domain the domain.
	 */
	public void setDomain(String domain) {
		this.domain = domain;
	}

	/**
	 * Accessor in reading on the account.
	 * @return the account.
	 */
	public String getAccount() {
		return account;
	}

	/**
	 * Accessor in writing on the account.
	 * @param account the account.
	 */
	public void setAccount(String account) {
		this.account = account;
	}

	
	
}
