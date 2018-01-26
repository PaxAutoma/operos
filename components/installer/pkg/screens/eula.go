/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package screens

import (
	"image"

	"github.com/paxautoma/operos/components/common/widgets"
)

func EULAScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	screen := widgets.NewScreen()
	screen.Title = "End-User License Agreement"
	screen.Message = "Please review and accept the agreement below."

	par := widgets.NewPar("eula", eula)
	par.Focusable = true
	par.Bounds = image.Rect(1, 1, 78, 15)
	screen.Content = par

	screen.OnNext = func() error {
		screenSet.Forward(1)
		return nil
	}

	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}

	return screen
}

const eula = `Software End User License Agreement
===================================

This End User License Agreement, including the Order Form which by this
reference is incorporated herein (this "Agreement"), is a binding agreement
between Pax Automa Systems Incorporated ("Pax Automa") and the person or
entity identified on the Order Form as the licensee of the Software
("Licensee").

PAX AUTOMA PROVIDES THE SOFTWARE SOLELY ON THE TERMS AND CONDITIONS SET
FORTH IN THIS AGREEMENT AND ON THE CONDITION THAT LICENSEE ACCEPTS AND
COMPLIES WITH THEM. BY CLICKING THE "ACCEPT" BUTTON OR BY CHECKING THE
"ACCEPT" BOX ON THE ORDER FORM YOU A. ACCEPT THIS AGREEMENT AND AGREE THAT
LICENSEE IS LEGALLY BOUND BY ITS TERMS; AND B. REPRESENT AND WARRANT THAT:
I. YOU ARE OF LEGAL AGE TO ENTER INTO A BINDING AGREEMENT; AND (II) IF
LICENSEE IS A CORPORATION, GOVERNMENTAL ORGANIZATION OR OTHER LEGAL ENTITY,
YOU HAVE THE RIGHT, POWER AND AUTHORITY TO ENTER INTO THIS AGREEMENT ON
BEHALF OF LICENSEE AND BIND LICENSEE TO ITS TERMS. IF LICENSEE DOES NOT
AGREE TO THE TERMS OF THIS AGREEMENT, PAX AUTOMA WILL NOT AND DOES NOT
LICENSE THE SOFTWARE TO LICENSEE AND YOU MUST NOT DOWNLOAD OR INSTALL THE
SOFTWARE OR DOCUMENTATION.

NOTWITHSTANDING ANYTHING TO THE CONTRARY IN THIS AGREEMENT OR YOUR OR
LICENSEE'S ACCEPTANCE OF THE TERMS AND CONDITIONS OF THIS AGREEMENT, NO
LICENSE IS GRANTED (WHETHER EXPRESSLY, BY IMPLICATION OR OTHERWISE) UNDER
THIS AGREEMENT, AND THIS AGREEMENT EXPRESSLY EXCLUDES ANY RIGHT, CONCERNING
ANY SOFTWARE THAT LICENSEE DID NOT ACQUIRE LAWFULLY OR THAT IS NOT A
LEGITIMATE, AUTHORIZED COPY OF PAX AUTOMA'S SOFTWARE.

1. Definitions. For purposes of this Agreement, the following terms have the
   following meanings:

    a. "Authorized Users" means those individuals authorized by Licensee to
       use the Software pursuant to the license granted under this
       Agreement, each of whom must accept the Authorized User Terms of Use
       attached as Annex 1 prior to using the Software.

    b. "Documentation" means user manuals, technical manuals and any other
       materials provided by Pax Automa, in printed, electronic or other
       form, that describe the installation, operation, use or technical
       specifications of the Software.

    c. "Intellectual Property Rights" means any and all registered and
       unregistered rights granted, applied for or otherwise now or
       hereafter in existence under or related to any patent, copyright,
       trade-mark, trade secret, database protection or other intellectual
       property rights laws, and all similar or equivalent rights or forms
       of protection, in any part of the world.

    d. "License Fees" means the license fees, including all taxes thereon,
       paid by Licensee for the license granted under this Agreement.

    e. "Order Form" means the order form filled out and submitted by or on
       behalf of Licensee, and accepted by Pax Automa, for Licensee's
       purchase of the license for the Software granted under this
       Agreement.  Most Order Forms are completed on Pax Automa's website at
       the time of download of the Software.

    f. "Person" means an individual, corporation, partnership, joint
       venture, governmental authority, unincorporated organization, trust,
       association or other entity.

    g. "Software" means Pax Automa's Operos software program for which
       Licensee is purchasing a license, and/or such other software programs
       as expressly set forth in the Order Form.

    h. "Term" has the meaning set forth in Section 11.1.

    i. "Third Party" means any Person other than Licensee or Pax Automa.

    j. "Update" has the meaning set forth in Section 7.

2. License Grant and Scope. Subject to and conditional on Licensee's payment
   of the License Fees, as applicable, and Licensee's strict compliance with
   all terms and conditions set forth in this Agreement, Pax Automa hereby
   grants to Licensee a non-exclusive, non-transferable, non-sublicensable,
   limited license during the Term to use, solely by and through its
   Authorized Users, the Software and Documentation, solely as set forth in
   this Section 2 and subject to all conditions and limitations set forth in
   Section 4 or elsewhere in this Agreement. This license grants Licensee
   the right, exercisable solely by and through Licensee's Authorized Users,
   to:

    2.1. Download and install in accordance with the Documentation one (1)
         copy of the Software on one (1) computer plus the number of other
         machines specified on the Order Form, each owned or leased, and
         controlled by, Licensee. In addition to the foregoing, Licensee has
         the right to make one (1) copy of the Software solely for archival
         and backup purposes. All copies of the Software made by Licensee:

        a. will be the exclusive property of Pax Automa;

        b. will be subject to the terms and conditions of this Agreement;
           and

        c. must include all trade-mark, copyright, patent, and other
           Intellectual Property Rights notices contained in the original.

    2.2. Use and run the Software as properly installed in accordance with
         this Agreement and the Documentation, solely as set forth in the
         Documentation and solely for Licensee's internal business purposes.

    2.3. Download or otherwise make one (1) copy of the Documentation per
         copy of the Software permitted to be downloaded and installed in
         accordance with this Agreement and use such Documentation solely in
         support of its licensed use of the Software in accordance herewith.
         All copies of the Documentation made by Licensee:

        a. will be the exclusive property of Pax Automa;

        b. will be subject to the terms and conditions of this Agreement;
           and

        c. must include all trade-mark, copyright, patent and other
           Intellectual Property Rights notices contained in the original.

    2.4. Transfer any copy of the Software from one computer to another,
         provided that the number of computers on which the Software is
         installed at any one time does not exceed the number permitted
         under Section 2.1,

3. Third-Party Materials. The Software may include software, content, data
   or other materials, including related documentation, that are owned by
   Persons other than Pax Automa and that are provided to Licensee on
   licensee terms that are in addition to and/or different from those
   contained in this Agreement ("Third-party Licenses"). A list of all
   materials, if any, included in the Software and provided under Third-
   party Licenses can be found at www.paxautoma.com/products/operos/credits
   and the applicable Third- party Licenses are accessible via links
   therefrom. Licensee is bound by and shall comply with all Third-party
   Licenses. Any breach by Licensee or any of its Authorized Users of any
   Third-party License is also a breach of this Agreement.

4. Use Restrictions.  Licensee shall not, and shall require its Authorized
   Users not to, directly or indirectly:

    a. use (including make any copies of) the Software or Documentation
       beyond the scope of the license granted under Section 2;

    b. provide any other Person, including any subcontractor, independent
       contractor, affiliate, or service provider of Licensee, with access
       to or use of the Software or Documentation;

    c. modify, translate, adapt or otherwise create derivative works or
       improvements, whether or not patentable, of the Software or
       Documentation or any part thereof;

    d. combine the Software or any part thereof with, or incorporate the
       Software or any part thereof in, any other programs;

    e. reverse engineer, disassemble, decompile, decode or otherwise attempt
       to derive or gain access to the source code of the Software or any
       part thereof;

    f. remove, delete, alter or obscure any trade-marks or any copyright,
       trade-mark, patent or other intellectual property or proprietary
       rights notices provided on or with the Software or Documentation,
       including any copy thereof;

    g. except as expressly set forth in Section 2.1 and Section 2.3, copy
       the Software or Documentation, in whole or in part;

    h. use the Software or Documentation in, or in association with, the
       design, construction, maintenance or operation of any hazardous
       environments or systems, including:

        i. power generation systems;

        ii. aircraft navigation or communication systems, air traffic
            control systems or any other transport management systems;

        iii. safety-critical applications, including medical or life-support
             systems, vehicle operation applications or any police, fire or
             other safety response systems; and

        iv. military or aerospace applications, weapons systems or
            environments;

    i. use the Software or Documentation in violation of any law, regulation
       or rule; or

    j. use the Software or Documentation for purposes of competitive
       analysis of the Software, the development of a competing software
       product or service or any other purpose that is to Pax Automa's
       commercial disadvantage.

5. Responsibility for Use of Software. Licensee is responsible and liable
   for all uses of the Software and Documentation through access thereto
   provided by Licensee, directly or indirectly. Specifically, and without
   limiting the generality of the foregoing, Licensee is responsible and
   liable for all actions and failures to take required actions with respect
   to the Software and Documentation by its Authorized Users or by any other
   Person to whom Licensee or an Authorized User may provide access to or
   use of the Software or Documentation, whether such access or use is
   permitted by or in violation of this Agreement.

6. Compliance Measures.

    6.1. The Software contains technological copy protection or other
         security features designed to prevent unauthorized use of the
         Software, including features to protect against any use of the
         Software that is prohibited under Section 4. Licensee shall not,
         and shall not attempt to, remove, disable, circumvent or otherwise
         create or implement any workaround to, any such copy protection or
         security features.

    6.2. During the Term, Pax Automa may, at any time and in Pax Automa's
         sole discretion, audit Licensee's use of the Software to ensure
         Licensee's compliance with this Agreement.  Pax Automa also may, in
         its sole discretion, audit Licensee's systems within three (3)
         months after the end of the Term to ensure Licensee has ceased use
         of the Software and removed all copies of the Software from such
         systems as required hereunder. Licensee shall fully cooperate with
         Pax Automa's personnel conducting such audits and provide all
         access requested by Pax Automa to records, systems, equipment,
         information and personnel, including machine IDs, serial numbers
         and related information.

    6.3. If any of the measures taken or implemented under this Section 6
         determines that Licensee's use of the Software exceeds or exceeded
         the use permitted by this Agreement, then:

        a. If License Fees were not required to be paid by Licensee for the
           license to the Software granted under this Agreement, such
           license and this Agreement may be terminated by Pax Automa
           without notice.

        b. If License Fees were required to be paid by Licensee for the
           license of the Software granted under this Agreement, Licensee
           shall, within five (5) days following the earlier to occur of (A)
           the date of such determination by Licensee or (B) Pax Automa's
           written notification thereof, pay to Pax Automa the retroactive
           License Fees for such excess use and obtain and pay for a valid
           license to bring Licensee's use into compliance with this
           Agreement. In determining Licensee Fee payable in accordance with
           the foregoing, i. unless Licensee can demonstrate otherwise by
           documentary evidence, all excess use of the Software shall be
           deemed to have commenced on the commencement date of this
           Agreement or, if later, the completion date of any audit
           previously conducted by Pax Automa hereunder and continued
           uninterrupted thereafter, and ii. the rates for such licenses
           shall be determined without regard to any discount to which
           Licensee may have been entitled had such use been properly
           licensed before its commencement (or deemed commencement). Pax
           Automa's remedies set forth in this Section 6.3 are cumulative
           and are in addition to, and not in lieu of, all other remedies
           Pax Automa may have at law or in equity, whether under this
           Agreement or otherwise.

7. Maintenance and Support. During the Term, Pax Automa will provide to
   Licensee such updates, upgrades, bug fixes, patches and other error
   corrections (collectively, "Updates") as Pax Automa makes generally
   available free of charge to all licensees of the Software.  Pax Automa
   may develop and provide Updates in its sole discretion, and Licensee
   agrees that Pax Automa has no obligation to develop any Updates at all or
   for particular issues. Licensee further agrees that all Updates will be
   deemed Software, and related documentation will be deemed Documentation,
   all subject to all terms and conditions of this Agreement. Licensee
   acknowledges that Pax Automa may provide Updates via download from a
   website designated by Pax Automa and that Licensee's receipt thereof will
   require an internet connection, which connection is Licensee's sole
   responsibility. Pax Automa has no obligation to provide Updates via any
   other media. Maintenance and support services do not include any new
   version or new release of the Software that Pax Automa may issue as a
   separate or new product, and Pax Automa may determine whether any
   issuance qualifies as a new version, new release or Update in its sole
   discretion.

8. Collection and Use of Information.

    8.1. Licensee acknowledges that Pax Automa may, directly or indirectly
         through the services of Third Parties, collect and store
         information regarding use of the Software and about equipment on
         which the Software is installed or through which it otherwise is
         accessed and used, through:

        a. the provision of maintenance and support services; and

        b. security measures included in the Software as described in
           Section 6.

    8.2. Licensee agrees that Pax Automa may use such information for any
         purpose related to any use of the Software by Licensee or on
         Licensee's equipment, including but not limited to:

        a. improving the performance of the Software or developing Updates;
           and

        b. verifying Licensee's compliance with the terms of this Agreement
           and enforcing Pax Automa's rights, including all Intellectual
           Property Rights in and to the Software.

9. Intellectual Property Rights. Licensee acknowledges and agrees that the
   Software and Documentation are provided under license, and not sold, to
   Licensee. Licensee does not acquire any ownership interest in the
   Software or Documentation under this Agreement, or any other rights
   thereto, other than to use the same in accordance with the license
   granted and subject to all terms, conditions and restrictions under this
   Agreement. Pax Automa and its licensors reserve and shall retain their
   entire right, title and interest in and to the Software and all
   Intellectual Property Rights arising out of or relating to the Software,
   except as expressly granted to Licensee in this Agreement. Licensee shall
   safeguard all Software (including all copies thereof) from infringement,
   misappropriation, theft, misuse or unauthorized access. Licensee shall
   promptly notify Pax Automa if Licensee becomes aware of any infringement
   of Pax Automa's Intellectual Property Rights in the Software and fully
   cooperate with Pax Automa in any legal action taken by Pax Automa to
   enforce its Intellectual Property Rights.

10. Payment. All License Fees are payable in the manner set forth in the
    Order Form and are non-refundable, except as expressly set forth herein.
    Any renewal of the license or maintenance and support services hereunder
    shall not be effective until the fees for such renewal have been paid in
    full.

11. Term and Termination.

    11.1. This Agreement and the license granted hereunder shall remain in
          effect for the term set forth on the Order Form (the "Initial
          Term") or until earlier terminated as set forth herein, and shall
          automatically renew for additional terms equal to the length of
          the Initial Term (together with the Initial Term, the "Term")
          unless either party provides thirty (30) days' written notice to
          the other party prior to such renewal that it does not intend to
          renew the Term.

    11.2. Licensee may terminate this Agreement by ceasing to use and
          destroying all copies of the Software and Documentation.

    11.3. Pax Automa may terminate this Agreement, effective upon written
          notice to Licensee, if Licensee breaches this Agreement and such
          breach: (i) is incapable of cure; or (ii) being capable of cure,
          remains uncured ten (10) days after Pax Automa provides written
          notice thereof.

    11.4. Pax Automa may terminate this Agreement, effective immediately, if
          Licensee files an assignment in bankruptcy or has a bankruptcy
          order made against it under any bankruptcy or insolvency law,
          makes or seeks to make a general assignment for the benefit of its
          creditors or applies for, or consents to, the appointment of a
          trustee, receiver, receiver-manager, monitor or custodian for all
          or a substantial part of its property.

    11.5. Upon expiration or earlier termination of this Agreement, the
          license granted hereunder shall also terminate, and Licensee shall
          cease using and destroy all copies of the Software and
          Documentation. No expiration or termination shall affect
          Licensee's obligation to pay all Licensee Fees that may have
          become due before such expiration or termination, or entitle
          Licensee to any refund.

12. Warranty Disclaimer. THE SOFTWARE AND DOCUMENTATION ARE PROVIDED TO USER
    "AS IS" AND WITH ALL FAULTS AND DEFECTS WITHOUT CONDITION OR WARRANTY OF
    ANY KIND. TO THE MAXIMUM EXTENT PERMITTED UNDER APPLICABLE LAW, PAX
    AUTOMA, ON ITS OWN BEHALF AND ON BEHALF OF ITS AFFILIATES AND ITS AND
    THEIR RESPECTIVE LICENSORS AND SERVICE PROVIDERS, EXPRESSLY DISCLAIMS
    ALL CONDITIONS AND WARRANTIES, WHETHER EXPRESS, IMPLIED, STATUTORY, OR
    OTHERWISE, WITH RESPECT TO THE SOFTWARE AND DOCUMENTATION, INCLUDING ALL
    IMPLIED CONDITIONS AND WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
    PARTICULAR PURPOSE, TITLE, QUIET POSSESSION AND NON-INFRINGEMENT, AND
    WARRANTIES THAT MAY ARISE OUT OF COURSE OF DEALING, COURSE OF
    PERFORMANCE, USAGE OR TRADE PRACTICE. WITHOUT LIMITING THE FOREGOING,
    PAX AUTOMA PROVIDES NO CONDITION, WARRANTY OR UNDERTAKING, AND MAKES NO
    REPRESENTATION OF ANY KIND THAT THE LICENSED SOFTWARE WILL MEET USER'S
    REQUIREMENTS, ACHIEVE ANY INTENDED RESULTS, BE COMPATIBLE OR WORK WITH
    ANY OTHER SOFTWARE, APPLICATIONS, SYSTEMS OR SERVICES, OPERATE WITHOUT
    INTERRUPTION, MEET ANY PERFORMANCE OR RELIABILITY STANDARDS OR BE ERROR
    FREE OR THAT ANY ERRORS OR DEFECTS CAN OR WILL BE CORRECTED.

13. Limitation of Liability. TO THE FULLEST EXTENT PERMITTED UNDER
    APPLICABLE LAW:

    13.1. IN NO EVENT WILL PAX AUTOMA OR ITS AFFILIATES, OR ANY OF ITS OR
          THEIR RESPECTIVE LICENSORS OR SERVICE PROVIDERS, BE LIABLE TO USER
          OR ANY THIRD PARTY FOR: (a) ANY: (i) USE, INTERRUPTION, DELAY OR
          INABILITY TO USE THE SOFTWARE; (ii) LOST REVENUES OR PROFITS;
          (iii) DELAYS, INTERRUPTION OR LOSS OF SERVICES, BUSINESS OR
          GOODWILL; (iv) LOSS OR CORRUPTION OF DATA; (v) LOSS RESULTING FROM
          SYSTEM OR SYSTEM SERVICE FAILURE, MALFUNCTION OR SHUTDOWN; (vi)
          FAILURE TO ACCURATELY TRANSFER, READ OR TRANSMIT INFORMATION;
          (vii) FAILURE TO UPDATE OR PROVIDE CORRECT INFORMATION; (viii)
          SYSTEM INCOMPATIBILITY OR PROVISION OF INCORRECT COMPATIBILITY
          INFORMATION; (ix) BREACHES IN SYSTEM SECURITY; OR (b) ANY
          CONSEQUENTIAL, INCIDENTAL, INDIRECT, SPECIAL, PUNITIVE OR
          EXEMPLARY DAMAGES, IN EACH CASE WHETHER ARISING OUT OF OR IN
          CONNECTION WITH THIS AGREEMENT, BREACH OF CONTRACT, TORT
          (INCLUDING NEGLIGENCE) OR OTHERWISE, REGARDLESS OF WHETHER SUCH
          DAMAGES WERE FORESEEABLE AND WHETHER OR NOT PAX AUTOMA WAS ADVISED
          OF THE POSSIBILITY OF SUCH DAMAGES.

    13.2. IN NO EVENT WILL THE COLLECTIVE AGGREGATE LIABILITY OF PAX AUTOMA
          AND ITS AFFILIATES, INCLUDING ANY OF ITS OR THEIR RESPECTIVE
          LICENSORS AND SERVICE PROVIDERS, UNDER OR IN CONNECTION WITH THIS
          AGREEMENT OR ITS SUBJECT MATTER, UNDER ANY LEGAL OR EQUITABLE
          THEORY, INCLUDING BREACH OF CONTRACT, TORT (INCLUDING NEGLIGENCE),
          STRICT LIABILITY AND OTHERWISE, EXCEED THE TOTAL AMOUNT PAID TO
          PAX AUTOMA UNDER THIS AGREEMENT FOR THE SOFTWARE THAT IS THE
          SUBJECT OF THE CLAIM.

    13.3. THE LIMITATIONS SET FORTH IN SECTION 13.1 AND SECTION 13.2 SHALL
          APPLY EVEN IF USER'S REMEDIES UNDER THIS AGREEMENT FAIL OF THEIR
          ESSENTIAL PURPOSE.

14. Export Regulation. The Software and Documentation may be subject to
    Canadian export control laws. Licensee shall not, directly or
    indirectly, export, re-export or release the Software or Documentation
    to, or make the Software or Documentation accessible from, any
    jurisdiction or country to which export, re-export or release is
    prohibited by law, rule or regulation. Licensee shall comply with all
    applicable federal laws, regulations and rules and complete all required
    undertakings (including obtaining any necessary export license or other
    governmental approval), before exporting, re-exporting, releasing or
    otherwise making the Software or Documentation available outside Canada.

15. Miscellaneous.

    15.1. This Agreement is governed by and construed in accordance with the
          laws of the Province of British Columbia and the federal laws of
          Canada applicable therein without regard to its conflict of law
          provisions. The United Nations Convention on Contracts for the
          International Sale of Goods does not apply to this Agreement.

    15.2. All disputes arising out of or in connection with this Agreement
          will be referred to and finally resolved by arbitration under the
          rules of the British Columbia International Commercial Arbitration
          Centre. The appointing authority will be the British Columbia
          International Commercial Arbitration Centre.  The case will be
          adjudicated by a single arbitrator and will be administered by the
          British Columbia International Commercial Arbitration Centre in
          accordance with its rules. The place of arbitration will be
          Vancouver, British Columbia, Canada.  The language of the
          arbitration will be English. Notwithstanding the foregoing, Pax
          Automa may seek and obtain injunctive relief in any jurisdiction
          in any court of competent jurisdiction and you agree that this
          Agreement is specifically enforceable by Pax Automa through
          injunctive relief and other equitable remedies without proof of
          monetary damages.

    15.3. Pax Automa will not be responsible or liable to Licensee, or
          deemed in default or breach hereunder by reason of any failure or
          delay in the performance of its obligations hereunder where such
          failure or delay is due to strikes, labour disputes, civil
          disturbances, riot, rebellion, invasion, epidemic, hostilities,
          war, terrorist attack, embargo, natural disaster, acts of God,
          flood, tsunami, fire, sabotage, fluctuations or non-availability
          of electrical power, heat, light, air conditioning or Licensee
          equipment, loss and destruction of property or any other
          circumstances or causes beyond Pax Automa's reasonable control.

    15.4. All notices, requests, consents, claims, demands, waivers and
          other communications hereunder shall be in writing and shall be
          deemed to have been given: (i) when delivered by hand (with
          written confirmation of receipt); (ii) when received by the
          addressee if sent by a nationally recognized overnight courier
          (receipt requested); (iii) on the date sent by e-mail (with
          confirmation of transmission) if sent during normal business hours
          of the recipient, and on the next business day if sent after
          normal business hours of the recipient; or (iv) on the third day
          after the date mailed, by certified or registered mail, return
          receipt requested, postage prepaid. Such communications must be
          sent to the respective parties at the addresses set forth on the
          Order Form (or to such other address as may be designated by a
          party from time to time in accordance with this Section 15.4).

    15.5. This Agreement, together with the Order Form and all other
          documents that are incorporated by reference herein, constitutes
          the sole and entire agreement between Licensee and Pax Automa with
          respect to the subject matter contained herein, and supersedes all
          prior and contemporaneous understandings, agreements,
          representations and warranties, both written and oral, with
          respect to such subject matter.

    15.6. Licensee shall not assign or otherwise transfer any of its rights,
          or delegate or otherwise transfer any of its obligations or
          performance, under this Agreement, in each case whether
          voluntarily, involuntarily, by operation of law or otherwise,
          without Pax Automa's prior written consent, which consent Pax
          Automa may give or withhold in its sole discretion. No delegation
          or other transfer will relieve Licensee of any of its obligations
          or performance under this Agreement. Any purported assignment,
          delegation or transfer in violation of this Section 15.6 is void.
          Pax Automa may freely assign or otherwise transfer all or any of
          its rights, or delegate or otherwise transfer all or any of its
          obligations or performance under this Agreement without Licensee's
          consent. This Agreement is binding upon and enures to the benefit
          of the parties hereto and their respective permitted successors
          and assigns.

    15.7. This Agreement is for the sole benefit of the parties hereto and
          their respective successors and permitted assigns and nothing
          herein, express or implied, is intended to or shall confer on any
          other Person any legal or equitable right, benefit or remedy of
          any nature whatsoever under or by reason of this Agreement.

    15.8. This Agreement may only be amended, modified or supplemented by an
          agreement in writing signed by each party hereto. No waiver by any
          party of any of the provisions hereof shall be effective unless
          explicitly set forth in writing and signed by the party so
          waiving. Except as otherwise set forth in this Agreement, no
          failure to exercise, or delay in exercising, any right, remedy,
          power or privilege arising from this Agreement shall operate or be
          construed as a waiver thereof; nor shall any single or partial
          exercise of any right, remedy, power or privilege hereunder
          preclude any other or further exercise thereof or the exercise of
          any other right, remedy, power or privilege.

    15.9. If any term or provision of this Agreement is invalid, illegal, or
          unenforceable in any jurisdiction, such invalidity, illegality or
          unenforceability shall not affect any other term or provision of
          this Agreement or invalidate or render unenforceable such term or
          provision in any other jurisdiction.

    15.10. For purposes of this Agreement, (a) the words "include,"
           "includes," and "including" shall be deemed to be followed by the
           words "without limitation"; (b) the word "or" is not exclusive;
           and (c) the words "herein," "hereof," "hereby," "hereto," and
           "hereunder" refer to this Agreement as a whole. Unless the
           context otherwise requires, references herein: (i) to Sections
           and Exhibits refer to the Sections of, and Exhibits attached to,
           this Agreement; (ii) to an agreement, instrument, or other
           document means such agreement, instrument, or other document as
           amended, supplemented, and modified from time to time to the
           extent permitted by the provisions thereof; and (iii) to a
           statute means such statute as amended from time to time and
           includes any successor legislation thereto and any regulations
           promulgated thereunder. This Agreement shall be construed without
           regard to any presumption or rule requiring construction or
           interpretation against the party drafting an instrument or
           causing any instrument to be drafted. The Order Form and all
           other documents that are incorporated by reference herein shall
           be construed with, and as an integral part of, this Agreement to
           the same extent as if they were set forth verbatim herein. Unless
           otherwise stated, all dollar amounts referred to in this
           Agreement are stated in Canadian dollars.

    15.11. The parties confirm that it is their express wish that this
           Agreement, as well as any other documents related to this
           Agreement, including notices, schedules and authorizations, have
           been and shall be drawn up in the English language only. Les
           parties aux présentes confirment leur volonté expresse que cette
           convention, de même que tous les documents s'y rattachant, y
           compris tous avis, annexes et autorisations s'y rattachant,
           soient rédigés en langue anglaise seulement.

    15.12. The headings in this Agreement are for reference only and do not
           affect the interpretation of this Agreement.
`
